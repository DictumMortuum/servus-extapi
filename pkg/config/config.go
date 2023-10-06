package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Databases map[string]string `config:"databases"`
	Port      string            `config:"port"`
}

var (
	Cfg Config
)

func Load() error {
	loader := confita.NewLoader(
		file.NewBackend("/etc/conf.d/servusrc.yml"),
	)

	err := loader.Load(context.Background(), &Cfg)
	if err != nil {
		return err
	}

	return nil
}
