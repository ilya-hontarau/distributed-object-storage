package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ilya-hontarau/distributed-object-storage/internal/gateway"
)

type Minio struct {
	client     *minio.Client
	bucketName string
}

type MinioConfig struct {
	Endpoint   string
	BucketName string
	AccessKey  string
	SecretKey  string
}

func NewMinio(ctx context.Context, cfg MinioConfig) (*Minio, error) {
	// TODO: refactor signature
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create a new client: %w", err)
	}
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket exists: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return &Minio{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

func (m *Minio) Upload(ctx context.Context, id string, file io.Reader, contentLength int) error {
	_, err := m.client.PutObject(ctx, m.bucketName, id, file, int64(contentLength), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}
	return nil
}

func (m *Minio) Download(ctx context.Context, id string) (io.Reader, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	_, err = object.Stat()
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.StatusCode == http.StatusNotFound {
			return nil, gateway.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get object stat: %w", err)
	}
	return object, nil
}
