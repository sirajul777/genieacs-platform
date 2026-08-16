package artifact

import (
	"context"
	"io"
)

type Store interface {
	Save(ctx context.Context, name string, r io.Reader) (string, error)
	Open(ctx context.Context, name string) (io.ReadCloser, error)
}
