package artifact

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestLocalStoreSaveAndOpen(t *testing.T) {
	s := NewLocalStore(t.TempDir())
	name, err := s.Save(context.Background(), "test.artifact", strings.NewReader("hello"))
	if err != nil { t.Fatal(err) }
	f, err := s.Open(context.Background(), name)
	if err != nil { t.Fatal(err) }
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil { t.Fatal(err) }
	if string(b) != "hello" { t.Fatalf("got %q", b) }
}
