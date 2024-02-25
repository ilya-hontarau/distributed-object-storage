package storage

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/minio/minio-go/v7"
)

type Minio struct {
	client     *minio.Client
	bucketName string
}

func NewMinio(ctx context.Context, endpoint, bucketName string) (*Minio, error) {
	client, err := minio.New(endpoint, nil)
	if err != nil {
		return nil, err
	}
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}
	return &Minio{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (m *Minio) Upload(ctx context.Context, id string, file fs.File) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}
	_, err = m.client.PutObject(ctx, m.bucketName, id, file, stat.Size(), minio.PutObjectOptions{})
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
	return object, nil
}
