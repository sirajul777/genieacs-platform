package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Default struct {
	Key   string
	Value any
}

type Loader struct {
	service string
	file    string
}

func NewLoader(service, file string) *Loader {
	return &Loader{service: service, file: file}
}

func (l *Loader) Load(target any, defaults ...Default) error {
	v := viper.New()
	v.SetEnvPrefix(l.service)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for _, item := range defaults {
		v.SetDefault(item.Key, item.Value)
	}

	if l.file != "" {
		v.SetConfigFile(l.file)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
	}

	return v.Unmarshal(target)
}
