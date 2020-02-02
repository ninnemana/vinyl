package datastore

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/ninnemana/vinyl/pkg/log"
	"github.com/ninnemana/vinyl/pkg/vinyl"
	"golang.org/x/sync/errgroup"

	"cloud.google.com/go/datastore"
	discogs "github.com/irlndts/go-discogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const (
	Entity = "vinyl"
)

var (
	ErrInvalidLogger = errors.New("the provided logger was not valid")
)

type Service struct {
	client      *datastore.Client
	discogs     *discogs.Client
	environment string
	log         *zap.Logger
}

func Register(server *grpc.Server) error {

	zlg, err := log.Init()
	if err != nil {
		return errors.Wrap(err, "failed to create logger")
	}

	svc, err := New(context.Background(), zlg, os.Getenv("GCE_PROJECT_ID"))
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
		Token:     "ChvDgMlrKNxsaMFyUISklJcyjCTwhxihcbOAMuCh",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create discogs client")
	}

	opts := []option.ClientOption{}
	client, err := datastore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create datastore client")
	}

	return &Service{
		discogs: disc,
		client:  client,
		log:     log,
	}, nil
}

// List retrieves all the entries that are associated with the user.
func (s *Service) List(p *vinyl.ListParams, srv vinyl.Vinyl_ListServer) error {

	q := datastore.NewQuery(Entity).Namespace(s.environment)
	if p.Artist != "" {
		q = q.Filter("artist =", p.Artist)
	}

	if p.Type != "" {
		q = q.Filter("type =", p.Type)
	}

	if p.Title != "" {
		q = q.Filter("title =", p.Title)
	}

	iter := s.client.Run(srv.Context(), q)

Loop:
	for {
		var res vinyl.ReleaseSource
		_, err := iter.Next(&res)
		switch err {
		case nil:
			if err := srv.Send(&res); err != nil {
				return errors.Wrap(err, "failed to send record of server")
			}
		case iterator.Done:
			break Loop
		default:
			return errors.Wrap(err, "failed to retrieve records from the datastore")
		}
	}

	return nil
}

func (s *Service) Get(ctx context.Context, p *vinyl.GetParams) (*vinyl.Result, error) {
	if p.GetId() == "" {
		return nil, vinyl.ErrInvalidGetParams
	}

	var (
		stored *vinyl.Result
		result *vinyl.Result
	)

	g, ctx := errgroup.WithContext(ctx)
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

		result = &vinyl.Result{
			// Catno:       res.Catno,
			// Format:      res.Formats,
			Id:          int64(res.ID),
			Title:       res.Title,
			ResourceUrl: res.ResourceURL,
			Thumb:       res.Thumb,
			// Year:        year,
			// Type:        res.,
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if stored != nil {
		return stored, nil
	}

	if result != nil {
		return result, nil
	}

	return nil, vinyl.ErrNotFound
}

func (s *Service) Search(p *vinyl.SearchParams, srv vinyl.Vinyl_SearchServer) error {
	s.log.Debug(
		"Searching for matching records against Discogs",
		zap.String("query", p.GetQuery()),
		zap.String("releaseTitle", p.GetReleaseTitle()),
		zap.String("type", p.GetType()),
		zap.String("title", p.GetTitle()),
		zap.String("credit", p.GetCredit()),
		zap.String("artist", p.GetArtist()),
		zap.String("anv", p.GetAnv()),
		zap.String("label", p.GetLabel()),
		zap.String("genre", p.GetGenre()),
		zap.String("country", p.GetCountry()),
		zap.String("format", p.GetFormat()),
		zap.String("contributor", p.GetContributor()),
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
		srv.Send(&vinyl.ReleaseSource{
			Catno:       res.Catno,
			Format:      strings.Join(res.Format, ","),
			Id:          int64(res.ID),
			Title:       res.Title,
			ResourceUrl: res.ResourceURL,
			Thumb:       res.Thumb,
			Year:        year,
			Type:        res.Type,
		})
	}
	return nil
}

func (s *Service) Store(ctx context.Context, p *vinyl.Result) (*vinyl.Result, error) {

	return nil, nil
}

func (s *Service) Health(_ context.Context, _ *vinyl.HealthRequest) (*vinyl.HealthResponse, error) {
	return &vinyl.HealthResponse{}, nil
}

func (s *Service) getStored(ctx context.Context, p *vinyl.GetParams) (*vinyl.Result, error) {
	q := datastore.NewQuery(Entity).Namespace(s.environment)
	q = q.Filter("ID =", p.GetId())

	var res vinyl.Result
	_, err := s.client.Run(ctx, q).Next(&res)
	switch err {
	case nil:
		return &res, nil
	case iterator.Done:
		return nil, vinyl.ErrNotFound
	default:
		return nil, errors.Wrap(err, "failed to retrieve record from the datastore")
	}
}
