package main

import (
	"LogAnalyse/app/service/file/rpc/pkg"
	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/kitex_gen/file"
	"LogAnalyse/app/shared/log"
	"LogAnalyse/app/shared/response"
	"context"
	"fmt"
	"time"
)

// FileServiceImpl implements the last service interface defined in the IDL.
type FileServiceImpl struct {
	Minio *pkg.MinioManager
}

// UploadFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) UploadFile(ctx context.Context, req *file.UploadFileReq) (resp *file.UploadFileResp, err error) {
	// TODO: Your code here...
	resp = new(file.UploadFileResp)

	bucketName := fmt.Sprintf("%v", req.Id)
	objName := time.Now().Format(time.RFC3339) + ".log"
	url, err := s.Minio.UploadFile(ctx, bucketName, objName)
	if err != nil {
		log.Zlogger.Error("generate upload file url failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGenerateUrl)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	resp.Url = url.String()
	return
}

// DownloadFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) DownloadFile(ctx context.Context, req *file.DownloadFileReq) (resp *file.DownloadFileResp, err error) {
	// TODO: Your code here...
	resp = new(file.DownloadFileResp)

	bucketName := fmt.Sprintf("%v", req.Id)
	url, err := s.Minio.DownLoadFile(ctx, bucketName, req.FileName)
	if err != nil {
		log.Zlogger.Error("generate download file url failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGenerateUrl)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	resp.Url = url.String()
	return
}

// ListFile implements the FileServiceImpl interface.
func (s *FileServiceImpl) ListFile(ctx context.Context, req *file.ListFileReq) (resp *file.ListFileResp, err error) {
	// TODO: Your code here...
	resp = new(file.ListFileResp)

	bucketName := fmt.Sprintf("%v", req.Id)
	infoCh, err := s.Minio.ListFile(ctx, bucketName)
	if err != nil {
		log.Zlogger.Error("get file list failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGetFileList)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	if len(infoCh) == 0 {
		return resp, nil
	}
	for info := range infoCh {
		resp.FileList = append(resp.FileList, &file.FileInfo{Name: info.Key, Size: info.Size, LastModifyTime: info.LastModified.UnixNano()})
	}
	return
}
