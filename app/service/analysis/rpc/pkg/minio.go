package pkg

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

type MinioManager struct {
	client *minio.Client
}

func NewMinioManager(cli *minio.Client) *MinioManager {
	return &MinioManager{client: cli}
}

func (m *MinioManager) UploadFile(ctx context.Context, path, jobName string, id int64) (string, error) {
	bucketName := fmt.Sprintf("%v", id)
	objectName := fmt.Sprintf("%soutput.json", jobName)
	_, err := m.client.FPutObject(ctx, bucketName, objectName, path, minio.PutObjectOptions{
		ContentType: "application/json",
	})
	return objectName, err
}
