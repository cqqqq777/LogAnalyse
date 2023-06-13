package pkg

import (
	"LogAnalyse/app/shared/kitex_gen/analyse"
	"LogAnalyse/app/shared/kitex_gen/analyse/analyseservice"
	"context"
)

type AnalyseManager struct {
	cli analyseservice.Client
}

func NewAnalysisManager(cli analyseservice.Client) AnalyseManager {
	return AnalyseManager{cli: cli}
}

func (a AnalyseManager) Analyse(ctx context.Context, jobID, id int64, url, field, jobName string) {
	a.cli.Analyse(ctx, &analyse.AnalyseReq{JobId: jobID, UserId: id, Url: url, Field: field, JobName: jobName})
}
