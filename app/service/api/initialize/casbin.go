package initialize

import (
	"LogAnalyse/app/service/api/config"
	"LogAnalyse/app/shared/log"
	"github.com/casbin/casbin/v2"
)

func InitCasbin() {
	e, err := casbin.NewEnforcer("./app/shared/middleware/model.conf", "./app/shared/middleware/policy.csv")
	if err != nil {
		log.Zlogger.Fatal("new casbin enforcer failed err:" + err.Error())
	}
	config.GlobalCasbinEnforcer = e
}
