package main

import (
	"context"
	"encoding/json"

	"LogAnalyse/app/service/analysis/internal"
	"LogAnalyse/app/service/analysis/rpc/pkg"
	"LogAnalyse/app/service/analysis/rpc/pkg/mr"
	"LogAnalyse/app/shared/kitex_gen/analyse"
	"LogAnalyse/app/shared/kitex_gen/common"
	"LogAnalyse/app/shared/log"
)

// AnalyseServiceImpl implements the last service interface defined in the IDL.
type AnalyseServiceImpl struct {
	Producer
	MinioManager
}

type Producer interface {
	Produce(msg pkg.Message) error
}

type MinioManager interface {
	UploadFile(ctx context.Context, path, jobName string, id int64) (string, error)
}

// Analyse implements the AnalyseServiceImpl interface.
func (s *AnalyseServiceImpl) Analyse(ctx context.Context, req *analyse.AnalyseReq) (resp *common.NilResp, err error) {
	// TODO: Your code here...
	paths, err := pkg.DownloadFile(req.Url, req.UserId)
	if err != nil {
		log.Zlogger.Warn("get file failed err:" + err.Error())
		err = s.Produce(pkg.Message{
			JobID:  req.JobId,
			Status: internal.StatusFailed,
		})
		if err != nil {
			log.Zlogger.Warn("publish nsq message failed err:" + err.Error())
		}
		return
	}

	// get data from MapReduce
	kvs, err := mr.StartWordCount(paths, req.Field)
	if err != nil {
		log.Zlogger.Warn("calculate failed err:" + err.Error())
		err = s.Produce(pkg.Message{
			JobID:  req.JobId,
			Status: internal.StatusFailed,
		})
		if err != nil {
			log.Zlogger.Warn("publish nsq message failed err:" + err.Error())
		}
		return
	}
	data, err := json.Marshal(kvs)
	if err != nil {
		log.Zlogger.Warn("marshal failed err:" + err.Error())
		err = s.Produce(pkg.Message{
			JobID:  req.JobId,
			Status: internal.StatusFailed,
		})
		if err != nil {
			log.Zlogger.Warn("publish nsq message failed err:" + err.Error())
		}
		return
	}

	// create output file and upload it to minio
	path := pkg.GetFilePath(req.UserId) + "output.json"
	err = pkg.CreateFile(path, data)
	if err != nil {
		log.Zlogger.Warn("create output file failed err:" + err.Error())
		err = s.Produce(pkg.Message{
			JobID:  req.JobId,
			Status: internal.StatusFailed,
		})
		if err != nil {
			log.Zlogger.Warn("publish nsq message failed err:" + err.Error())
		}
		return
	}

	outputPath, err := s.UploadFile(ctx, req.JobName, path, req.UserId)
	if err != nil {
		log.Zlogger.Warn("upload output file failed err:" + err.Error())
		err = s.Produce(pkg.Message{
			JobID:  req.JobId,
			Status: internal.StatusFailed,
		})
		if err != nil {
			log.Zlogger.Warn("publish nsq message failed err:" + err.Error())
		}
		return
	}

	pkg.DeleteFile(paths)
	err = s.Produce(pkg.Message{
		JobID:          req.JobId,
		Status:         internal.StatusSuccess,
		ConsequentFile: outputPath,
	})
	if err != nil {
		log.Zlogger.Warn("publish nsq message failed err:" + err.Error())
	}
	return
}
