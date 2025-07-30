package pkg

import (
	"fmt"
	"os"
)

func GenerateUrl(filename string) string {
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	bucket := os.Getenv("MINIO_BUCKETS")
	if minioEndpoint != "localhost:9000" {
		return fmt.Sprintf("https://%s/%s/%s", minioEndpoint, bucket, filename)
	}
	return fmt.Sprintf("http://%s/%s/%s", minioEndpoint, bucket, filename)
}
