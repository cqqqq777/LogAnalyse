// Code generated by Kitex v0.5.1. DO NOT EDIT.

package fileservice

import (
	file "LogAnalyse/app/shared/kitex_gen/file"
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return fileServiceServiceInfo
}

var fileServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "FileService"
	handlerType := (*file.FileService)(nil)
	methods := map[string]kitex.MethodInfo{
		"UploadFile":   kitex.NewMethodInfo(uploadFileHandler, newFileServiceUploadFileArgs, newFileServiceUploadFileResult, false),
		"DownloadFile": kitex.NewMethodInfo(downloadFileHandler, newFileServiceDownloadFileArgs, newFileServiceDownloadFileResult, false),
		"ListFile":     kitex.NewMethodInfo(listFileHandler, newFileServiceListFileArgs, newFileServiceListFileResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "file",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.5.1",
		Extra:           extra,
	}
	return svcInfo
}

func uploadFileHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*file.FileServiceUploadFileArgs)
	realResult := result.(*file.FileServiceUploadFileResult)
	success, err := handler.(file.FileService).UploadFile(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFileServiceUploadFileArgs() interface{} {
	return file.NewFileServiceUploadFileArgs()
}

func newFileServiceUploadFileResult() interface{} {
	return file.NewFileServiceUploadFileResult()
}

func downloadFileHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*file.FileServiceDownloadFileArgs)
	realResult := result.(*file.FileServiceDownloadFileResult)
	success, err := handler.(file.FileService).DownloadFile(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFileServiceDownloadFileArgs() interface{} {
	return file.NewFileServiceDownloadFileArgs()
}

func newFileServiceDownloadFileResult() interface{} {
	return file.NewFileServiceDownloadFileResult()
}

func listFileHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*file.FileServiceListFileArgs)
	realResult := result.(*file.FileServiceListFileResult)
	success, err := handler.(file.FileService).ListFile(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newFileServiceListFileArgs() interface{} {
	return file.NewFileServiceListFileArgs()
}

func newFileServiceListFileResult() interface{} {
	return file.NewFileServiceListFileResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) UploadFile(ctx context.Context, req *file.UploadFileReq) (r *file.UploadFileResp, err error) {
	var _args file.FileServiceUploadFileArgs
	_args.Req = req
	var _result file.FileServiceUploadFileResult
	if err = p.c.Call(ctx, "UploadFile", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DownloadFile(ctx context.Context, req *file.DownloadFileReq) (r *file.DownloadFileResp, err error) {
	var _args file.FileServiceDownloadFileArgs
	_args.Req = req
	var _result file.FileServiceDownloadFileResult
	if err = p.c.Call(ctx, "DownloadFile", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListFile(ctx context.Context, req *file.ListFileReq) (r *file.ListFileResp, err error) {
	var _args file.FileServiceListFileArgs
	_args.Req = req
	var _result file.FileServiceListFileResult
	if err = p.c.Call(ctx, "ListFile", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
