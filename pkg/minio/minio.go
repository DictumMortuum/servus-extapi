package minio

import (
	"github.com/DictumMortuum/servus-extapi/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Init() (*minio.Client, error) {
	return minio.New(config.Cfg.Minio.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Cfg.Minio.Key, config.Cfg.Minio.Secret, ""),
		Secure: config.Cfg.Minio.SSL,
	})
}
