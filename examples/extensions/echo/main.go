package main

import (
	"context"
	"log"

	"github.com/sirajul777/genieacs-platform/internal/agent/extension"
)

type Echo struct{}

func (Echo) Name() string { return "echo" }
func (Echo) Execute(_ context.Context, req extension.Request) (extension.Response, error) {
	return extension.Response{Payload: req.Payload}, nil
}

func main() {
	registry := extension.NewRegistry(Echo{})
	if _, ok := registry.Get("echo"); !ok {
		log.Fatal("extension not registered")
	}
}
