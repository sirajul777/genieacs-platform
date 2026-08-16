package artifact

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type LocalStore struct{ Root string }

func NewLocalStore(root string) *LocalStore { return &LocalStore{Root: root} }

func (s *LocalStore) path(name string) string { return filepath.Join(s.Root, filepath.Base(name)) }

func (s *LocalStore) Save(ctx context.Context, name string, r io.Reader) (string, error) {
	if err := ctx.Err(); err != nil { return "", err }
	if err := os.MkdirAll(s.Root, 0o750); err != nil { return "", err }
	path := s.path(name)
	f, err := os.Create(path)
	if err != nil { return "", err }
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil { return "", err }
	return path, nil
}

func (s *LocalStore) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil { return nil, err }
	return os.Open(s.path(name))
}
