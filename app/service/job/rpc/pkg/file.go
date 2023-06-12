package pkg

import (
	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/kitex_gen/file"
	"LogAnalyse/app/shared/kitex_gen/file/fileservice"
	"context"
	"fmt"
)

type FileManager struct {
	client fileservice.Client
}

func NewFileManager(client fileservice.Client) FileManager {
	return FileManager{client: client}
}

func (f FileManager) DownloadFile(ctx context.Context, id int64, fileName string) (string, error) {
	resp, err := f.client.DownloadFile(ctx, &file.DownloadFileReq{Id: id, FileName: fileName})
	if err != nil {
		return "", err
	}
	if resp.BaseResp.Code != errz.Success {
		return "", fmt.Errorf(resp.BaseResp.Msg)
	}
	return resp.Url, nil
}
