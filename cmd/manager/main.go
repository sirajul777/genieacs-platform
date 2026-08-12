package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	managerconfig "github.com/sirajul777/genieacs-platform/internal/manager/config"
	"github.com/sirajul777/genieacs-platform/internal/manager/server"
	"github.com/sirajul777/genieacs-platform/internal/shared/database"
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
		Use:     "manager",
		Short:   "Run the GenieACS platform manager service",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := managerconfig.Load(configFile)
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

			db, err := database.Connect(ctx, cfg.Database)
			if err != nil {
				return err
			}
			defer db.Close()

			if err := server.New(cfg.Server, log, db).Start(ctx); err != nil {
				log.Error("manager server stopped", zap.Error(err))
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "path to config file")
	return cmd
}
