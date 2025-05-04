package attachmanager

import (
	"bytes"
	"context"
	"errors"
	"log"

	"github.com/lim-bo/calendar/backend/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type AttachManager struct {
	cli    *minio.Client
	bucket string
}

type MinioCfg struct {
	Address    string
	User       string
	Pass       string
	BucketName string
}

func New(cfg MinioCfg) *AttachManager {
	client, err := minio.New(cfg.Address, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Pass, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &AttachManager{
		cli:    client,
		bucket: cfg.BucketName,
	}
}

func (am *AttachManager) Load(file *models.FileLoad) error {
	_, err := am.cli.PutObject(context.Background(), am.bucket, file.Name, bytes.NewReader(file.Data), int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return errors.New("loading file error: " + err.Error())
	}
	return nil
}
