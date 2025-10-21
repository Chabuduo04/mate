package services

import (
	"context"
	"io"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"github.com/Chabuduo04/mate/back_end/config"
)

type StorageService interface {
	Upload(file io.Reader, key string, fileName string) error
}

type KodoService struct {
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewKodoService() *KodoService {
	cfg := config.GetConfig()
	return &KodoService{
		AccessKey: cfg.Kodo.AccessKey,
		SecretKey: cfg.Kodo.SecretKey,
		Bucket:    cfg.Kodo.Bucket,
	}
}

func (s *KodoService) Upload(file io.Reader, key string, fileName string) error {
	mac := credentials.NewCredentials(s.AccessKey, s.SecretKey)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err := uploadManager.UploadReader(context.Background(), file, &uploader.ObjectOptions{
		BucketName: s.Bucket,
		ObjectName: &key,
		FileName:   fileName,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
