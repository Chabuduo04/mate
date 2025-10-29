package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"github.com/Chabuduo04/mate/back_end/config"
)

type StorageService interface {
	Upload(file io.Reader, key string, fileName string) error
	UploadAudio(file multipart.File, fileName string) (string, error)
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

func (s *KodoService) UploadAudio(file multipart.File, fileName string) (string, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(fileName)
	if fileExt == "" {
		fileExt = ".wav" // 默认扩展名
	}
	fmt.Println(fileExt)
	key := "voice-chat/" + uuid.New().String() + fileExt

	uploadReader := bytes.NewReader(fileBytes)
	err = s.Upload(uploadReader, key, fileName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s/%s", "http://", config.GetConfig().Kodo.Host, key), nil
}
