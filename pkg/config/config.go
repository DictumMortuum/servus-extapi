package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type ModemConfig struct {
	Host  string            `config:"host"`
	Modem string            `config:"modem"`
	User  string            `config:"user"`
	Pass  string            `config:"pass"`
	Voip  string            `config:"voip"`
	Extra map[string]string `config:"extra"`
}

type Config struct {
	Databases map[string]string      `config:"databases"`
	Port      string                 `config:"port"`
	Modem     map[string]ModemConfig `config:"modem"`
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
