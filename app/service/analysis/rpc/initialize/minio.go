package initialize

import (
	"LogAnalyse/app/service/analysis/rpc/config"
	"LogAnalyse/app/shared/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
	endpoint := config.GlobalServerConfig.MinioInfo.Endpoint
	accessKeyID := config.GlobalServerConfig.MinioInfo.AccessKeyID
	secretAccessKey := config.GlobalServerConfig.MinioInfo.SecretAccessKey

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Zlogger.Fatal("init minio client failed err:" + err.Error())
	}

	return minioClient
}
