package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"github.com/heetch/confita/backend/flags"
)

type ModemConfig struct {
	Host  string            `config:"host"`
	Modem string            `config:"modem"`
	User  string            `config:"user"`
	Pass  string            `config:"pass"`
	Voip  string            `config:"voip"`
	Extra map[string]string `config:"extra"`
}

type DecoConfig struct {
	Host   string `config:"host"`
	Pass   string `config:"pass"`
	Folder string `config:"folder"`
}

type MinioConfig struct {
	Host   string `config:"minio_host"`
	Key    string `config:"minio_key"`
	Secret string `config:"minio_secret"`
	SSL    bool   `config:"minio_ssl"`
}

type ApiFlashConfig struct {
	Host string `config:"apiflash_host"`
	Key  string `config:"apiflash_key"`
}

type Config struct {
	Databases        map[string]string      `config:"databases"`
	Port             string                 `config:"port"`
	FileExporterPort string                 `config:"file"`
	Modem            map[string]ModemConfig `config:"modem"`
	Deco             DecoConfig             `config:"deco"`
	Minio            MinioConfig            `config:"minio"`
	ApiFlash         ApiFlashConfig         `config:"apiflash"`
	ServusHost       string                 `config:"servus_host"`
	ServusUser       string                 `config:"servus_user"`
	ServusPass       string                 `config:"servus_pass"`
}

var (
	Cfg Config
)

func Load() error {
	loader := confita.NewLoader(
		flags.NewBackend(),
		env.NewBackend(),
		file.NewBackend("/etc/conf.d/servusrc.yml"),
	)

	Cfg = Config{
		FileExporterPort: ":10005",
	}

	err := loader.Load(context.Background(), &Cfg)
	if err != nil {
		return err
	}

	return nil
}
