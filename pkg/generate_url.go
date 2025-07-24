package pkg

import (
	"fmt"
	"os"
)

func GenerateUrl(filename string) string {
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	bucket := os.Getenv("MINIO_BUCKETS")
	return fmt.Sprintf("%s/%s/%s", minioEndpoint, bucket, filename)
}
