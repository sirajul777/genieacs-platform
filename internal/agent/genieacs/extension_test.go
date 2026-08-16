package genieacs

import (
	"context"
	"testing"

	"github.com/sirajul777/genieacs-platform/internal/agent/extension"
)

func TestExtensionExecuteUsesCommand(t *testing.T) {
	ext := New("printf")
	out, err := ext.Execute(context.Background(), extension.Request{Payload: []byte(`{"action":"hello","args":{"x":"world"}}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Payload) == 0 {
		t.Fatal("expected response payload")
	}
}
