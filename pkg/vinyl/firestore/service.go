package firestore

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ninnemana/vinyl/pkg/auth"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/firestore"
	discogs "github.com/ninnemana/go-discogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Entity = "vinyls"
)

var (
	ErrInvalidLogger = errors.New("the provided logger was not valid")
)

type Service struct {
	client        *firestore.Client
	discogs       discogs.Discogs
	log           *zap.Logger
	tokenizer     auth.Tokenizer
	initTimestamp time.Time
	hostname      string
}

type Config struct {
	Logger          *zap.Logger
	GoogleProjectID string
	DiscogsAPIKey   string
	Hostname        string
	Tokenizer       auth.Tokenizer
}

func New(ctx context.Context, cfg Config) (*Service, error) {
	if cfg.Logger == nil {
		return nil, ErrInvalidLogger
	}

	disc, err := discogs.New(&discogs.Options{
		UserAgent: "Some Agent",
		Token:     cfg.DiscogsAPIKey,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discogs client")
	}

	client, err := firestore.NewClient(ctx, cfg.GoogleProjectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	return &Service{
		discogs:       disc,
		client:        client,
		log:           cfg.Logger,
		hostname:      cfg.Hostname,
		tokenizer:     cfg.Tokenizer,
		initTimestamp: time.Now().UTC(),
	}, nil
}

// List retrieves all the entries that are associated with the user.
func (s *Service) List(ctx context.Context, p *vinyl.ListParams) (*vinyl.ListResponse, error) {
	q := s.client.Collection(Entity).OrderBy("year", firestore.Desc)

	if p.Artist != "" {
		q = q.Where("artist", "==", p.GetArtist())
	}

	if p.Type != "" {
		q = q.Where("type", "==", p.GetType())
	}

	if p.Title != "" {
		q = q.Where("title", "==", p.GetTitle())
	}

	results := []*vinyl.Release{}
	getter := func(iter *firestore.DocumentIterator) error {
		doc, err := iter.Next()
		switch err {
		case nil:
			var res vinyl.Release
			if err := doc.DataTo(&res); err != nil {
				return errors.Wrap(err, "document was not valid type")
			}

			results = append(results, &res)

			return nil
		case iterator.Done:
			return iterator.Done
		default:
			return errors.Wrap(err, "failed to retrieve records from the firestore")
		}
	}

	it := q.Documents(ctx)

	for {
		err := getter(it)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	return &vinyl.ListResponse{
		Results: results,
	}, nil
}

func (s *Service) Middleware() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		s.tokenizer.Authenticator,
	}
}

func (s *Service) Get(ctx context.Context, p *vinyl.GetParams) (*vinyl.Master, error) {
	if p.GetId() == "" {
		return nil, vinyl.ErrInvalidGetParams
	}

	var (
		stored *vinyl.Master
		result *vinyl.Master
	)

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		r, err := s.getStored(ctx, p)
		if err == vinyl.ErrNotFound {
			return nil
		}

		if err != nil {
			return err
		}

		stored = r
		return nil
	})

	g.Go(func() error {
		id, err := strconv.Atoi(p.GetId())
		if err != nil {
			return errors.New("failed to parse identifier `" + p.GetId() + "`")
		}

		s.log.Debug("querying discogs for release", zap.Int("releaseID", id))

		res, err := s.discogs.Master(id)
		if err != nil {
			return err
		}

		result = toMaster(res)

		return nil
	})

	if err := g.Wait(); err != nil {
		s.log.Error("failed to retrieve release", zap.Error(err))
		return nil, err
	}

	if stored != nil {
		s.log.Debug("found stored release")
		return stored, nil
	}

	if result != nil {
		s.log.Debug("found release on Discogs")
		return result, nil
	}

	s.log.Debug("no release found")

	return nil, vinyl.ErrNotFound
}

func (s *Service) Register(rpc *grpc.Server) error {
	s.log.Debug(
		"register RPC service",
		zap.Any("info", rpc.GetServiceInfo()),
	)

	// go func() {
	// 	ctx := context.Background()
	// 	for i := 0; i < 5; i++ {
	// 		conn, err := grpc.DialContext(
	// 			ctx,
	// 			"localhost:8080",
	// 			// grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
	// 			// 	Certificates:       certs,
	// 			// 	InsecureSkipVerify: true,
	// 			// })),
	// 			grpc.WithInsecure(),
	// 			grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	// 			grpc.WithUnaryInterceptor(drudge.UnaryClientInterceptor("vinyl")),
	// 			grpc.WithStreamInterceptor(drudge.StreamClientInterceptor("vinyl")),
	// 		)
	// 		if err != nil {
	// 			s.log.Error("client can't dial", zap.Error(err))
	// 			time.Sleep(time.Second * 1)
	// 			continue
	// 		}
	// 		defer conn.Close()

	// 		s.rpcClient = vinyl.NewVinylClient(conn)
	// 		if _, err := s.rpcClient.Health(ctx, &vinyl.HealthRequest{}); err != nil {
	// 			s.log.Error("health check is not responding", zap.Error(err))
	// 			time.Sleep(time.Second * 1)
	// 			continue
	// 		}

	// 		s.log.Debug("client connection established")
	// 		break
	// 	}
	// }()

	vinyl.RegisterVinylServer(rpc, s)

	return nil
}

func (s *Service) Route() string {
	return "/" + Entity
}

func (s *Service) Search(ctx context.Context, p *vinyl.SearchParams) (*vinyl.SearchResponse, error) {
	s.log.Debug(
		"Searching for matching records against Discogs",
		zap.Any("params", p),
	)

	search, err := s.discogs.Search(discogs.SearchRequest{
		Q:            p.GetQuery(),
		ReleaseTitle: p.GetReleaseTitle(),
		Type:         "master",
		Title:        p.GetTitle(),
		Credit:       p.GetCredit(),
		Artist:       p.GetArtist(),
		Anv:          p.GetAnv(),
		Label:        p.GetLabel(),
		Genre:        p.GetGenre(),
		Country:      p.GetCountry(),
		Format:       p.GetFormat(),
		Contributor:  p.GetContributor(),
	})
	if err != nil {
		s.log.Error(
			"failed to execute search operation",
			zap.Error(err),
		)
		return nil, errors.Wrap(err, "failed to execute search operation")
	}

	results := []*vinyl.ReleaseResponse{}
	for _, res := range search.Results {
		year, _ := strconv.ParseInt(res.Year, 0, 64)

		results = append(results, &vinyl.ReleaseResponse{
			Release: &vinyl.ReleaseSource{
				Catno:       res.Catno,
				Format:      strings.Join(res.Format, ","),
				Id:          int32(res.ID),
				Title:       res.Title,
				ResourceUrl: res.ResourceURL,
				Thumb:       res.Thumb,
				Year:        int32(year),
				Type:        res.Type,
			},
			Pagination: &vinyl.Pagination{
				PerPage: int32(search.Pagination.PerPage),
				Page:    int32(search.Pagination.Page),
				Items:   int32(search.Pagination.Items),
			},
		})
	}

	return &vinyl.SearchResponse{
		Results: results,
	}, nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := mux.CurrentRoute(r)
	if route == nil {
		s.log.Error("no route found", zap.String("path", r.URL.Path))
		return
	}

	sub := route.Subrouter()
	sub.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var params vinyl.SearchParams
		switch r.Method {
		case http.MethodGet:
			params = vinyl.SearchParams{
				Query:  r.URL.Query().Get("query"),
				Type:   r.URL.Query().Get("type"),
				Title:  r.URL.Query().Get("title"),
				Artist: r.URL.Query().Get("artist"),
			}
		case http.MethodPost:
			defer r.Body.Close()

			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "method not allow", http.StatusMethodNotAllowed)
			return
		}

		result, err := s.Search(r.Context(), &params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	sub.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		resp, err := s.Health(r.Context(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	sub.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		release, err := s.Get(r.Context(), &vinyl.GetParams{
			Id: mux.Vars(r)["id"],
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(release); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
	})
	sub.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.List(r.Context(), &vinyl.ListParams{
			Artist: r.URL.Query().Get("artist"),
			Type:   r.URL.Query().Get("type"),
			Title:  r.URL.Query().Get("title"),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	sub.ServeHTTP(w, r)
}

func (s *Service) Store(ctx context.Context, p *vinyl.Release) (*vinyl.Release, error) {
	_, err := s.client.Collection(Entity).Doc(
		strconv.Itoa(int(p.GetId())),
	).Set(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set release into the store")
	}

	return p, nil
}

func (s *Service) Health(_ context.Context, _ *vinyl.HealthRequest) (*vinyl.HealthResponse, error) {
	return &vinyl.HealthResponse{
		Uptime:  time.Since(s.initTimestamp).String(),
		Machine: s.hostname,
	}, nil
}

func (s *Service) getStored(ctx context.Context, p *vinyl.GetParams) (*vinyl.Master, error) {
	doc, err := s.client.Collection(Entity).Doc(p.GetId()).Get(ctx)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok {
			return nil, err
		}

		if s.Code() == codes.NotFound {
			return nil, vinyl.ErrNotFound
		}

		return nil, s.Err()
	}

	var res vinyl.Master
	if err := doc.DataTo(&res); err != nil {
		return nil, errors.Wrap(err, "failed to read document")
	}

	return &res, nil
}

func toRelease(res *discogs.Release) *vinyl.Release {
	result := &vinyl.Release{
		Id:          int32(res.ID),
		Title:       res.Title,
		ArtistsSort: res.ArtistsSort,
		DataQuality: res.DataQuality,
		Thumb:       res.Thumb,
		Community: &vinyl.Community{
			DataQuality: res.Community.DataQuality,
			Have:        int32(res.Community.Have),
			Rating: &vinyl.Rating{
				Average: res.Community.Rating.Average,
				Count:   int32(res.Community.Rating.Count),
			},
			Status: res.Community.Status,
			Submitter: &vinyl.Submitter{
				ResourceUrl: res.Community.Submitter.ResourceURL,
				Username:    res.Community.Submitter.Username,
			},
			Want: int32(res.Community.Want),
		},
		Country:           res.Country,
		DateAdded:         res.DateAdded,
		DateChanged:       res.DateChanged,
		EstimatedWeight:   int32(res.EstimatedWeight),
		Format:            nil,
		Genres:            res.Genres,
		LowestPrice:       float32(res.LowestPrice),
		MasterId:          int32(res.MasterID),
		MasterUrl:         res.MasterURL,
		Notes:             res.Notes,
		NumberForSale:     int32(res.NumForSale),
		Released:          res.Released,
		ReleasedFormatted: res.ReleasedFormatted,
		ResourceUrl:       res.ResourceURL,
		Series: func(series []discogs.Series) []string {
			r := make([]string, len(series))
			for i := range series {
				r[i] = series[i].Name
			}
			return r
		}(res.Series),
		Status: res.Status,
		Styles: res.Styles,
		Uri:    res.URI,
		Year:   int32(res.Year),
	}

	result.Artists = make([]*vinyl.ArtistSource, len(res.Artists))
	for i, artist := range res.Artists {
		result.Artists[i] = toArtist(artist)
	}

	result.ExtraArtists = make([]*vinyl.ArtistSource, len(res.ExtraArtists))
	for i, artist := range res.ExtraArtists {
		result.ExtraArtists[i] = toArtist(artist)
	}

	result.TrackList = make([]*vinyl.Track, len(res.Tracklist))
	for i, track := range res.Tracklist {
		result.TrackList[i] = &vinyl.Track{
			Duration: track.Duration,
			Position: track.Position,
			Title:    track.Title,
			Type:     track.Type,
		}

		result.TrackList[i].Artists = make([]*vinyl.ArtistSource, len(track.Artists))
		for j, artist := range track.Artists {
			result.TrackList[i].Artists[j] = toArtist(artist)
		}

		result.TrackList[i].Extraartists = make([]*vinyl.ArtistSource, len(track.Extraartists))
		for j, artist := range track.Extraartists {
			result.TrackList[i].Extraartists[j] = toArtist(artist)
		}
	}

	result.Companies = make([]*vinyl.Company, len(res.Companies))
	for i, cmp := range res.Companies {
		result.Companies[i] = &vinyl.Company{
			Catno:          cmp.Catno,
			EntityType:     cmp.EntityType,
			EntityTypeName: cmp.EntityTypeName,
			Id:             int32(cmp.ID),
			Name:           cmp.Name,
			ResourceUrl:    cmp.ResourceURL,
		}
	}

	result.Videos = make([]*vinyl.Video, len(res.Videos))
	for i, vid := range res.Videos {
		result.Videos[i] = &vinyl.Video{
			Description: vid.Description,
			Duration:    int32(vid.Duration),
			Embed:       vid.Embed,
			Title:       vid.Title,
			Uri:         vid.URI,
		}
	}

	result.Identifiers = make([]*vinyl.Identifier, len(res.Identifiers))
	for i, id := range res.Identifiers {
		result.Identifiers[i] = &vinyl.Identifier{
			Type:  id.Type,
			Value: id.Value,
		}
	}

	result.Images = make([]*vinyl.Image, len(res.Images))
	for i, img := range res.Images {
		result.Images[i] = &vinyl.Image{
			Height:      int32(img.Height),
			Width:       int32(img.Width),
			ResourceUrl: img.ResourceURL,
			Type:        img.Type,
			Uri:         img.URI,
			Uri150:      img.URI150,
		}
	}

	result.Labels = make([]*vinyl.LabelSource, len(res.Labels))
	for i, lbl := range res.Labels {
		result.Labels[i] = &vinyl.LabelSource{
			Catno:          lbl.Catno,
			EntityType:     lbl.EntityType,
			EntityTypeName: lbl.EntityTypeName,
			Id:             int32(lbl.ID),
			Name:           lbl.Name,
			ResourceUrl:    lbl.ResourceURL,
		}
	}

	result.Community.Contributors = make([]*vinyl.Contributor, len(res.Community.Contributors))
	for i, contrib := range res.Community.Contributors {
		result.Community.Contributors[i] = &vinyl.Contributor{
			ResourceUrl: contrib.ResourceURL,
			Username:    contrib.Username,
		}
	}

	return result
}

func toMaster(res *discogs.Master) *vinyl.Master {
	result := &vinyl.Master{
		Id:            int32(res.ID),
		Title:         res.Title,
		DataQuality:   res.DataQuality,
		Genres:        res.Genres,
		LowestPrice:   float32(res.LowestPrice),
		Notes:         res.Notes,
		NumberForSale: int32(res.NumForSale),
		ResourceUrl:   res.ResourceURL,
		Styles:        res.Styles,
		Uri:           res.URI,
		Year:          int32(res.Year),
	}

	result.Artists = make([]*vinyl.ArtistSource, len(res.Artists))
	for i, artist := range res.Artists {
		result.Artists[i] = toArtist(artist)
	}

	result.TrackList = make([]*vinyl.Track, len(res.Tracklist))
	for i, track := range res.Tracklist {
		result.TrackList[i] = &vinyl.Track{
			Duration: track.Duration,
			Position: track.Position,
			Title:    track.Title,
			Type:     track.Type,
		}

		result.TrackList[i].Artists = make([]*vinyl.ArtistSource, len(track.Artists))
		for j, artist := range track.Artists {
			result.TrackList[i].Artists[j] = toArtist(artist)
		}

		result.TrackList[i].Extraartists = make([]*vinyl.ArtistSource, len(track.Extraartists))
		for j, artist := range track.Extraartists {
			result.TrackList[i].Extraartists[j] = toArtist(artist)
		}
	}

	result.Videos = make([]*vinyl.Video, len(res.Videos))
	for i, vid := range res.Videos {
		result.Videos[i] = &vinyl.Video{
			Description: vid.Description,
			Duration:    int32(vid.Duration),
			Embed:       vid.Embed,
			Title:       vid.Title,
			Uri:         vid.URI,
		}
	}

	result.Images = make([]*vinyl.Image, len(res.Images))
	for i, img := range res.Images {
		result.Images[i] = &vinyl.Image{
			Height:      int32(img.Height),
			Width:       int32(img.Width),
			ResourceUrl: img.ResourceURL,
			Type:        img.Type,
			Uri:         img.URI,
			Uri150:      img.URI150,
		}
	}

	return result
}

func toArtist(artist discogs.ArtistSource) *vinyl.ArtistSource {
	return &vinyl.ArtistSource{
		Anv:         artist.Anv,
		Id:          int32(artist.ID),
		Join:        artist.Join,
		Name:        artist.Name,
		ResourceUrl: artist.ResourceURL,
		Role:        artist.Role,
		Tracks:      artist.Tracks,
	}
}

func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
