package config

import (
	"LogAnalyse/app/shared/kitex_gen/file/fileservice"
	"LogAnalyse/app/shared/kitex_gen/job/jobservice"
	"LogAnalyse/app/shared/kitex_gen/user/userservice"
	"github.com/casbin/casbin/v2"
)

var (
	GlobalServerConfig = &ServerConfig{}
	GlobalNacosConfig  = &NacosConfig{}

	GlobalUserClient userservice.Client
	GlobalFileClient fileservice.Client
	GlobalJobClient  jobservice.Client

	GlobalCasbinEnforcer *casbin.Enforcer
)
