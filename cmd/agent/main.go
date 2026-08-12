package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	agentconfig "github.com/sirajul777/genieacs-platform/internal/agent/config"
	"github.com/sirajul777/genieacs-platform/internal/agent/server"
	"github.com/sirajul777/genieacs-platform/internal/shared/logger"
	"github.com/sirajul777/genieacs-platform/internal/shared/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:     "agent",
		Short:   "Run the GenieACS platform agent service",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := agentconfig.Load(configFile)
			if err != nil {
				return err
			}
			log, err := logger.New(cfg.Logger)
			if err != nil {
				return err
			}
			defer func() { _ = log.Sync() }()

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			if err := server.New(cfg.Server, log).Start(ctx); err != nil {
				log.Error("agent server stopped", zap.Error(err))
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "path to config file")
	return cmd
}
