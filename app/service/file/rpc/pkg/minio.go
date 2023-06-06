package pkg

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

const expireTime = 1800

var reqParams url.Values

type MinioManager struct {
	client *minio.Client
}

func NewMinioManager(client *minio.Client) *MinioManager {
	return &MinioManager{client: client}
}

func (m *MinioManager) UploadFile(ctx context.Context, bucketName, objectName string) (*url.URL, error) {
	ok, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "cn-north-1", ObjectLocking: false})
		if err != nil {
			return nil, err
		}
	}
	return m.client.PresignedPutObject(ctx, bucketName, objectName, time.Second*expireTime)
}

func (m *MinioManager) DownLoadFile(ctx context.Context, bucketName, objectName string) (*url.URL, error) {
	ok, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "cn-north-1", ObjectLocking: false})
		if err != nil {
			return nil, err
		}
	}
	return m.client.PresignedGetObject(ctx, bucketName, objectName, time.Second*expireTime, reqParams)
}

func (m *MinioManager) ListFile(ctx context.Context, bucketName string) (<-chan minio.ObjectInfo, error) {
	ok, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "cn-north-1", ObjectLocking: false})
		if err != nil {
			return nil, err
		}
	}
	return m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true}), nil
}
