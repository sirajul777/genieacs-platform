package config

import (
	sharedconfig "github.com/sirajul777/genieacs-platform/internal/shared/config"
	"github.com/sirajul777/genieacs-platform/internal/shared/database"
	"github.com/sirajul777/genieacs-platform/internal/shared/logger"
)

type Config struct {
	Server   ServerConfig    `mapstructure:"server"`
	Logger   logger.Config   `mapstructure:"logger"`
	Database database.Config `mapstructure:"database"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

func Load(file string) (Config, error) {
	var cfg Config
	err := sharedconfig.NewLoader("manager", file).Load(&cfg,
		sharedconfig.Default{Key: "server.address", Value: ":8080"},
		sharedconfig.Default{Key: "logger.environment", Value: "development"},
		sharedconfig.Default{Key: "logger.level", Value: "info"},
		sharedconfig.Default{Key: "database.host", Value: "localhost"},
		sharedconfig.Default{Key: "database.port", Value: 5432},
		sharedconfig.Default{Key: "database.database", Value: "genieacs_platform"},
		sharedconfig.Default{Key: "database.user", Value: "genieacs"},
		sharedconfig.Default{Key: "database.password", Value: "genieacs"},
		sharedconfig.Default{Key: "database.ssl_mode", Value: "disable"},
	)
	return cfg, err
}
