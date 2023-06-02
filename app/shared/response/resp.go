package response

import (
	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/kitex_gen/common"
)

func NewBaseResp(err *errz.ErrZ) *common.BaseResp {
	if err == nil {
		return NewBaseResp(errz.NewErrZ())
	}
	return &common.BaseResp{
		Code: err.GetCode(),
		Msg:  err.Error(),
	}
}
