package middleware

import (
	"context"

	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/log"

	"github.com/casbin/casbin/v2"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type CasManager struct {
	enforcer *casbin.Enforcer
}

func newCasManager(e *casbin.Enforcer) *CasManager {
	return &CasManager{enforcer: e}
}

func CasbinAuth(e *casbin.Enforcer) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		casManager := newCasManager(e)
		role, ok := c.Get("accountIdentity")
		if !ok {
			role = "tourist"
		}
		ok, err := casManager.enforcer.Enforce(role, string(c.Path()), string(c.Method()))
		if err != nil {
			log.Zlogger.Warn("get authority failed err:" + err.Error())
			c.JSON(consts.StatusOK, utils.H{
				"code": errz.CodeServiceBusy,
				"msg":  "get authority failed",
			})
			return
		}
		if !ok {
			c.JSON(consts.StatusOK, utils.H{
				"code": errz.CodeNoPermission,
				"msg":  "no permission",
			})
			return
		}
		c.Next(ctx)
	}
}
