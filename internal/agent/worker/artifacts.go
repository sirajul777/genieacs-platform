package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/sirajul777/genieacs-platform/internal/agent/artifact"
)

type ArtifactWorker struct {
	TypeName string
	Command  string
	Store    artifact.Store
}

func NewArtifactWorker(typeName, command string, store artifact.Store) *ArtifactWorker {
	return &ArtifactWorker{TypeName: typeName, Command: command, Store: store}
}

func (w *ArtifactWorker) Type() string { return w.TypeName }

func (w *ArtifactWorker) Execute(ctx context.Context, job Job, progress ProgressReporter) (Result, error) {
	if w.Store == nil { return Result{}, errors.New("artifact store is required") }
	if w.Command == "" { return Result{}, errors.New("artifact command is required") }
	if err := progress.Report(ctx, 10, "artifact job started"); err != nil { return Result{}, err }

	switch w.TypeName {
	case "genieacs.backup":
		return w.backup(ctx, job, progress)
	case "genieacs.restore":
		return w.restore(ctx, job, progress)
	default:
		return Result{}, fmt.Errorf("unsupported artifact worker: %s", w.TypeName)
	}
}

func (w *ArtifactWorker) backup(ctx context.Context, job Job, progress ProgressReporter) (Result, error) {
	if err := progress.Report(ctx, 30, "creating backup"); err != nil { return Result{}, err }
	cmd := exec.CommandContext(ctx, w.Command, "backup")
	stdout, err := cmd.StdoutPipe()
	if err != nil { return Result{}, err }
	if err := cmd.Start(); err != nil { return Result{}, err }
	name := job.ID + ".backup"
	path, saveErr := w.Store.Save(ctx, name, stdout)
	waitErr := cmd.Wait()
	if saveErr != nil { return Result{}, saveErr }
	if waitErr != nil { return Result{}, fmt.Errorf("backup command: %w", waitErr) }
	if err := progress.Report(ctx, 100, "backup completed"); err != nil { return Result{}, err }
	return Result{Payload: []byte(path)}, nil
}

func (w *ArtifactWorker) restore(ctx context.Context, job Job, progress ProgressReporter) (Result, error) {
	if err := progress.Report(ctx, 30, "restoring backup"); err != nil { return Result{}, err }
	reader, err := w.Store.Open(ctx, string(job.Payload))
	if err != nil { return Result{}, err }
	defer reader.Close()
	cmd := exec.CommandContext(ctx, w.Command, "restore")
	stdin, err := cmd.StdinPipe()
	if err != nil { return Result{}, err }
	if err := cmd.Start(); err != nil { return Result{}, err }
	_, copyErr := io.Copy(stdin, reader)
	_ = stdin.Close()
	waitErr := cmd.Wait()
	if copyErr != nil { return Result{}, copyErr }
	if waitErr != nil { return Result{}, fmt.Errorf("restore command: %w", waitErr) }
	if err := progress.Report(ctx, 100, "restore completed"); err != nil { return Result{}, err }
	return Result{Payload: []byte("restored")}, nil
}
