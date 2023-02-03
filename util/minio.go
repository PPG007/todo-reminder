package util

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"time"
)

var MinioClient = &minioClient{}

type minioClient struct {
	c                  *minio.Client
	bucket             string
	urlExpiredDuration time.Duration
}

func init() {
	MinioClient.bucket = viper.GetString("minio.bucket")
	MinioClient.urlExpiredDuration = time.Duration(viper.GetInt("minio.urlExpiredSeconds")) * time.Second
	client, err := minio.New(viper.GetString("minio.endpoint"), &minio.Options{
		Creds: credentials.NewStaticV4(
			viper.GetString("minio.ak"),
			viper.GetString("minio.sk"),
			""),
	})
	if err != nil {
		panic(err)
	}
	MinioClient.c = client
}

func (m *minioClient) GetPreSignPutObjectUrl(ctx context.Context, fileName string) (string, error) {
	url, err := m.c.PresignedPutObject(ctx, m.bucket, fileName, m.urlExpiredDuration)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *minioClient) SignObjectUrl(ctx context.Context, fileName string) (string, error) {
	url, err := m.c.PresignedGetObject(ctx, m.bucket, fileName, m.urlExpiredDuration, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
