package logger

import "go.uber.org/zap"

type Config struct {
	Environment string `mapstructure:"environment"`
	Level       string `mapstructure:"level"`
}

func New(cfg Config) (*zap.Logger, error) {
	var zapCfg zap.Config
	if cfg.Environment == "production" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}
	if cfg.Level != "" {
		level := zap.NewAtomicLevel()
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			return nil, err
		}
		zapCfg.Level = level
	}
	return zapCfg.Build()
}
