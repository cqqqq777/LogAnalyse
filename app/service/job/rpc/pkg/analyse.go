package pkg

import "context"

type AnalyseManager struct {
}

func NewAnalysisManager(err error) AnalyseManager {
	return AnalyseManager{}
}

func (a AnalyseManager) Analyse(ctx context.Context, id int64, url, field string) {

}
