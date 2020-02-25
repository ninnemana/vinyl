package firestore

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/firestore"
	discogs "github.com/irlndts/go-discogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Entity = "vinyl"
)

var (
	ErrInvalidLogger = errors.New("the provided logger was not valid")
)

type Service struct {
	client      *firestore.Client
	discogs     *discogs.Client
	environment string
	log         *zap.Logger
}

func Register(server *grpc.Server) error {

	zlg, err := log.Init()
	if err != nil {
		return errors.Wrap(err, "failed to create logger")
	}

	svc, err := New(context.Background(), zlg, os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		return errors.Wrap(err, "failed to create micropost service")
	}

	vinyl.RegisterVinylServer(server, svc)

	return nil
}

func New(ctx context.Context, log *zap.Logger, projectID string) (*Service, error) {
	if log == nil {
		return nil, ErrInvalidLogger
	}

	disc, err := discogs.NewClient(&discogs.Options{
		UserAgent: "Some Agent",
		Token:     os.Getenv("DISCOGS_API_KEY"),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discogs client")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	return &Service{
		discogs: disc,
		client:  client,
		log:     log,
	}, nil
}

// List retrieves all the entries that are associated with the user.
func (s *Service) List(p *vinyl.ListParams, srv vinyl.Vinyl_ListServer) error {
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

	getter := func(iter *firestore.DocumentIterator) error {
		doc, err := iter.Next()
		switch err {
		case nil:
			var res vinyl.Release
			if err := doc.DataTo(&res); err != nil {
				return errors.Wrap(err, "document was not valid type")
			}

			if err := srv.Send(&res); err != nil {
				return errors.Wrap(err, "failed to send record of server")
			}

			return nil
		case iterator.Done:
			return iterator.Done
		default:
			return errors.Wrap(err, "failed to retrieve records from the firestore")
		}
	}

	it := q.Documents(srv.Context())

	for {
		err := getter(it)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Get(ctx context.Context, p *vinyl.GetParams) (*vinyl.Release, error) {
	if p.GetId() == "" {
		return nil, vinyl.ErrInvalidGetParams
	}

	var (
		stored *vinyl.Release
		result *vinyl.Release
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

		res, err := s.discogs.Database.Release(id)
		if err != nil {
			return err
		}

		result = toRelease(res)

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

func (s *Service) Search(p *vinyl.SearchParams, srv vinyl.Vinyl_SearchServer) error {
	s.log.Debug(
		"Searching for matching records against Discogs",
		zap.Any("params", p),
	)

	search, err := s.discogs.Search.Search(discogs.SearchRequest{
		Q:            p.GetQuery(),
		ReleaseTitle: p.GetReleaseTitle(),
		Type:         p.GetType(),
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
		return errors.Wrap(err, "failed to execute search operation")
	}

	for _, res := range search.Results {
		year, _ := strconv.ParseInt(res.Year, 0, 64)

		if err := srv.Send(&vinyl.ReleaseResponse{
			Release: &vinyl.ReleaseSource{
				Catno:       res.Catno,
				Format:      strings.Join(res.Format, ","),
				Id:          int64(res.ID),
				Title:       res.Title,
				ResourceUrl: res.ResourceURL,
				Thumb:       res.Thumb,
				Year:        year,
				Type:        res.Type,
			},
			Pagination: &vinyl.Pagination{
				PerPage: int64(search.Pagination.PerPage),
				Page:    int64(search.Pagination.Page),
				Items:   int64(search.Pagination.Items),
			},
		}); err != nil {
			return err
		}
	}

	return nil
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
	return &vinyl.HealthResponse{}, nil
}

func (s *Service) getStored(ctx context.Context, p *vinyl.GetParams) (*vinyl.Release, error) {
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

	var res vinyl.Release
	if err := doc.DataTo(&res); err != nil {
		return nil, errors.Wrap(err, "failed to read document")
	}

	return &res, nil
}

func toRelease(res *discogs.Release) *vinyl.Release {
	result := &vinyl.Release{
		Id:          int64(res.ID),
		Title:       res.Title,
		ArtistsSort: res.ArtistsSort,
		DataQuality: res.DataQuality,
		Thumb:       res.Thumb,
		Community: &vinyl.Community{
			DataQuality: res.Community.DataQuality,
			Have:        int64(res.Community.Have),
			Rating: &vinyl.Rating{
				Average: res.Community.Rating.Average,
				Count:   int64(res.Community.Rating.Count),
			},
			Status: res.Community.Status,
			Submitter: &vinyl.Submitter{
				ResourceUrl: res.Community.Submitter.ResourceURL,
				Username:    res.Community.Submitter.Username,
			},
			Want: int64(res.Community.Want),
		},
		Country:           res.Country,
		DateAdded:         res.DateAdded,
		DateChanged:       res.DateChanged,
		EstimatedWeight:   int64(res.EstimatedWeight),
		Format:            nil,
		Genres:            res.Genres,
		LowestPrice:       float32(res.LowestPrice),
		MasterId:          int64(res.MasterID),
		MasterUrl:         res.MasterURL,
		Notes:             res.Notes,
		NumberForSale:     int64(res.NumForSale),
		Released:          res.Released,
		ReleasedFormatted: res.ReleasedFormatted,
		ResourceUrl:       res.ResourceURL,
		Series:            res.Series,
		Status:            res.Status,
		Styles:            res.Styles,
		Uri:               res.URI,
		Year:              int64(res.Year),
	}

	result.Artists = make([]*vinyl.ArtistSource, len(res.Artists))
	for i, artist := range res.Artists {
		result.Artists[i] = toArtist(artist)
		fmt.Println(result.Artists[i])
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
			Id:             int64(cmp.ID),
			Name:           cmp.Name,
			ResourceUrl:    cmp.ResourceURL,
		}
	}

	result.Videos = make([]*vinyl.Video, len(res.Videos))
	for i, vid := range res.Videos {
		result.Videos[i] = &vinyl.Video{
			Description: vid.Description,
			Duration:    int64(vid.Duration),
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
			Height:      int64(img.Height),
			Width:       int64(img.Width),
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
			Id:             int64(lbl.ID),
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

func toArtist(artist discogs.ArtistSource) *vinyl.ArtistSource {
	return &vinyl.ArtistSource{
		Anv:         artist.Anv,
		Id:          int64(artist.ID),
		Join:        artist.Join,
		Name:        artist.Name,
		ResourceUrl: artist.ResourceURL,
		Role:        artist.Role,
		Tracks:      artist.Tracks,
	}
}
