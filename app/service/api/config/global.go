package config

import (
	"LogAnalyse/app/shared/kitex_gen/file/fileservice"
	"LogAnalyse/app/shared/kitex_gen/user/userservice"
)

var (
	GlobalServerConfig = &ServerConfig{}
	GlobalNacosConfig  = &NacosConfig{}

	GlobalUserClient userservice.Client
	GlobalFileClient fileservice.Client
)
