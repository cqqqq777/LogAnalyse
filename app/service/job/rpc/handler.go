package main

import (
	"LogAnalyse/app/service/job/internal"
	"LogAnalyse/app/service/job/rpc/model"
	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/log"
	"LogAnalyse/app/shared/response"
	"context"
	"github.com/bwmarrin/snowflake"
	"time"

	"LogAnalyse/app/service/job/rpc/dao"
	"LogAnalyse/app/shared/kitex_gen/job"
)

// JobServiceImpl implements the last service interface defined in the IDL.
type JobServiceImpl struct {
	AnalysisManager
	FileManager

	Dao *dao.Job
}

type AnalysisManager interface {
	Analyse(ctx context.Context, id int64, url, field string)
}

type FileManager interface {
	DownloadFile(ctx context.Context, id int64, fileName string) (string, error)
}

// CreateJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) CreateJob(ctx context.Context, req *job.CreateJobReq) (resp *job.CreateJobResp, err error) {
	// TODO: Your code here...
	resp = new(job.CreateJobResp)

	url, err := s.DownloadFile(ctx, req.UserId, req.FileName)
	if err != nil {
		log.Zlogger.Warn("get file from file service failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrDownFile)
		return resp, nil
	}

	sf, err := snowflake.NewNode(internal.JobSnowflakeNode)
	if err != nil {
		log.Zlogger.Warn("generate job id failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGenerateJobID)
		return resp, nil
	}
	jobID := sf.Generate().Int64()

	jobInfo := &model.Job{
		JobID:    jobID,
		UserID:   req.UserId,
		Status:   internal.StatusWait,
		FileName: req.FileName,
		JobName:  req.JobName,
	}

	err = s.Dao.CreateJob(jobInfo)
	if err != nil {
		log.Zlogger.Warn("create job in mysql failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrCreateJob)
		return resp, nil
	}

	go s.Analyse(ctx, req.UserId, url, req.Field)

	resp.BaseResp = response.NewBaseResp(nil)
	resp.JobInfo.JobId = jobID
	resp.JobInfo.UserId = req.UserId
	resp.JobInfo.JobName = req.JobName
	resp.JobInfo.Status = internal.StatusWait
	resp.JobInfo.FileName = req.FileName
	resp.JobInfo.CreateTime = time.Now().UnixNano()

	return resp, nil
}

// ListJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) ListJob(ctx context.Context, req *job.ListJobReq) (resp *job.ListJobResp, err error) {
	// TODO: Your code here...
	resp = new(job.ListJobResp)

	list, err := s.Dao.ListJob(req.UserId)
	if err != nil {
		log.Zlogger.Warn("query jobs failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGetJobList)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	for _, v := range list {
		resp.JobList = append(resp.JobList, &job.JobBaseInfo{JobId: v.JobID, CreateTime: v.CreateTime.UnixNano(), JobName: v.JobName})
	}

	return resp, nil
}

// GetJobInfo implements the JobServiceImpl interface.
func (s *JobServiceImpl) GetJobInfo(ctx context.Context, req *job.GetJobInfoReq) (resp *job.GetJobInfoResp, err error) {
	// TODO: Your code here...
	resp = new(job.GetJobInfoResp)

	info, err := s.Dao.GetJobInfo(req.JobId)
	if err != nil {
		log.Zlogger.Warn("get job info failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGetJobInfo)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	resp.JobInfo.JobId = info.JobID
	resp.JobInfo.UserId = info.UserID
	resp.JobInfo.Status = info.Status
	resp.JobInfo.JobName = info.JobName
	resp.JobInfo.FileName = info.FileName
	resp.JobInfo.ConsequentFile = info.ConsequentFile
	resp.JobInfo.CreateTime = info.CreateTime.UnixNano()

	return resp, nil
}

// UpdateJob implements the JobServiceImpl interface.
func (s *JobServiceImpl) UpdateJob(ctx context.Context, req *job.UpdateJobStatusReq) (resp *job.UpdateJobStatusResp, err error) {
	// TODO: Your code here...
	resp = new(job.UpdateJobStatusResp)

	err = s.Dao.UpdateJob(req.JobId, req.Status, req.ConsequentFile)
	if err != nil {
		log.Zlogger.Warn("update job failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrUpdateJob)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	return resp, nil
}
