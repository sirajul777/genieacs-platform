package config

import (
	sharedconfig "github.com/sirajul777/genieacs-platform/internal/shared/config"
	"github.com/sirajul777/genieacs-platform/internal/shared/logger"
)

type Config struct {
	Server ServerConfig  `mapstructure:"server"`
	Logger logger.Config `mapstructure:"logger"`
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
	)
	return cfg, err
}
