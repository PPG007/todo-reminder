package util

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"io"
	"time"
)

var MinioClient = &minioClient{}

type minioClient struct {
	client             *minio.Client
	bucket             string
	urlExpiredDuration time.Duration
}

func init() {
	if !viper.GetBool("minio.enabled") {
		return
	}
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
	MinioClient.client = client
	err = MinioClient.ensureBucket(context.Background())
	if err != nil {
		panic(err)
	}
}

func (m *minioClient) GetPreSignPutObjectUrl(ctx context.Context, fileName string) (string, error) {
	url, err := m.client.PresignedPutObject(ctx, m.bucket, fileName, m.urlExpiredDuration)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *minioClient) SignObjectUrl(ctx context.Context, fileName string) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucket, fileName, m.urlExpiredDuration, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (m *minioClient) ensureBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return err
	}
	if !exists {
		err := m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{
			ObjectLocking: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *minioClient) IsObjectExist(ctx context.Context, objectName string) bool {
	info, err := m.client.StatObject(ctx, m.bucket, objectName, minio.StatObjectOptions{})
	if err != nil || info.Err != nil {
		return false
	}
	return true
}

func (m *minioClient) PutObject(ctx context.Context, objectName string, object io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	size, err := buf.ReadFrom(object)
	if err != nil {
		return "", err
	}
	_, err = m.client.PutObject(ctx, m.bucket, objectName, buf, size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return m.SignObjectUrl(ctx, objectName)
}
