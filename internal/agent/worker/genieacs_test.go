package worker

import (
	"context"
	"testing"

	"github.com/sirajul777/genieacs-platform/internal/agent/extension"
)

type progressSink struct{ last int }
func (p *progressSink) Report(_ context.Context, progress int, _ string) error { p.last = progress; return nil }

type testExt struct{}
func (testExt) Name() string { return "genieacs" }
func (testExt) Execute(_ context.Context, req extension.Request) (extension.Response, error) {
	return extension.Response{Payload: req.Payload}, nil
}

func TestGenieACSWorkerExecutesInstall(t *testing.T) {
	r := extension.NewRunner(extension.NewRegistry(testExt{}))
	w := NewGenieACSWorker(r, "genieacs.install")
	p := &progressSink{}
	out, err := w.Execute(context.Background(), Job{Type: "genieacs.install", Payload: []byte(`{"args":{"version":"latest"}}`)}, p)
	if err != nil { t.Fatal(err) }
	if p.last != 100 { t.Fatalf("progress=%d", p.last) }
	if len(out.Payload) == 0 { t.Fatal("expected result") }
}

func TestGenieACSWorkerRejectsUnknownType(t *testing.T) {
	r := extension.NewRunner(extension.NewRegistry(testExt{}))
	w := NewGenieACSWorker(r, "genieacs.other")
	_, err := w.Execute(context.Background(), Job{Type: "genieacs.other"}, &progressSink{})
	if err == nil { t.Fatal("expected error") }
}
