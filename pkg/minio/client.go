package minio

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"store/pkg/config"
)

func NewClient(conf *config.GlobalConfig) (*minio.Client, error) {
	return minio.New(fmt.Sprintf("%s:%s", conf.Minio.Endpoint, conf.Minio.Port), &minio.Options{
		Creds: credentials.NewStaticV4(conf.Minio.AccessKey, conf.Minio.SecretKey, ""),
	})
}
