package vinyl

import "github.com/pkg/errors"

var (
	ErrNotFound         = errors.Errorf("no vinyl was found")
	ErrInvalidGetParams = errors.New("invalid parameters supplied when attempting to retrieve a vinyl")
)
