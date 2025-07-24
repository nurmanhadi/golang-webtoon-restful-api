package s3

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
)

type S3Storage interface {
	UploadFile(file *os.File, filename string) error
	RemoveFile(filename string) error
}
type s3Storage struct {
	ctx context.Context
	s3  *minio.Client
}

func NewS3Storage(ctx context.Context, s3 *minio.Client) S3Storage {
	return &s3Storage{
		ctx: ctx,
		s3:  s3,
	}
}
func (s *s3Storage) UploadFile(file *os.File, filename string) error {
	src, err := os.Open(file.Name())
	if err != nil {
		return err
	}
	defer src.Close()

	fileStat, err := src.Stat()
	if err != nil {
		return err
	}
	bucket := os.Getenv("MINIO_BUCKETS")
	_, err = s.s3.PutObject(s.ctx, bucket, filename, src, fileStat.Size(), minio.PutObjectOptions{ContentType: "image/webp"})
	if err != nil {
		return err
	}
	return nil
}
func (s *s3Storage) RemoveFile(filename string) error {
	bucket := os.Getenv("MINIO_BUCKETS")
	if err := s.s3.RemoveObject(s.ctx, bucket, filename, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}); err != nil {
		return err
	}
	return nil
}
