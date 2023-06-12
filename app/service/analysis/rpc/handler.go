package main

import (
	"LogAnalyse/app/service/analysis/rpc/pkg"
	"LogAnalyse/app/service/analysis/rpc/pkg/mr"
	analyse "LogAnalyse/app/shared/kitex_gen/analyse"
	common "LogAnalyse/app/shared/kitex_gen/common"
	"LogAnalyse/app/shared/log"
	"context"
	"encoding/json"
)

// AnalyseServiceImpl implements the last service interface defined in the IDL.
type AnalyseServiceImpl struct {
	Producer
}

type Producer interface {
	Produce(ctx context.Context, msg pkg.Message)
}

// Analyse implements the AnalyseServiceImpl interface.
func (s *AnalyseServiceImpl) Analyse(ctx context.Context, req *analyse.AnalyseReq) (resp *common.NilResp, err error) {
	// TODO: Your code here...
	paths, err := pkg.DownloadFile(req.Url, req.UserId)
	if err != nil {
		log.Zlogger.Warn("get file failed err:" + err.Error())
		s.Produce(ctx, pkg.Message{})
		return
	}

	kvs, err := mr.StartWordCount(paths, req.Field)
	if err != nil {
		log.Zlogger.Warn("calculate failed err:" + err.Error())
		s.Produce(ctx, pkg.Message{})
		return
	}
	data, err := json.Marshal(kvs)
	if err != nil {
		log.Zlogger.Warn("marshal failed err:" + err.Error())
		s.Produce(ctx, pkg.Message{})
		return
	}
	path := pkg.GetFilePath(req.UserId) + "output"
	err = pkg.CreateFile(path, data)
	if err != nil {
		log.Zlogger.Warn("create output file failed err:" + err.Error())
		s.Produce(ctx, pkg.Message{})
		return
	}

	err = pkg.UploadFile(path, req.UserId)
	if err != nil {
		log.Zlogger.Warn("upload output file failed err:" + err.Error())
		s.Produce(ctx, pkg.Message{})
		return
	}

	pkg.DeleteFile(paths)
	s.Produce(ctx, pkg.Message{})
	return
}
