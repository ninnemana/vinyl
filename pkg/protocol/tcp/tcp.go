package tcp

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ninnemana/vinyl/pkg/log"
)

const (
	tcpAddr   = "vinyltap:8000"
	projectID = "vinyl-registry"
)

func Serve() error {

	logger, err := log.Init()
	if err != nil {
		return errors.Wrap(err, "Failed to create logger")
	}

	server, err := NewServer(
		context.Background(),
		projectID,
		tcpAddr,
		logger,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create server")
	}
	defer server.Flush()

	return errors.Wrap(server.ListenAndServe(), "fell out of serving traffic")
}
