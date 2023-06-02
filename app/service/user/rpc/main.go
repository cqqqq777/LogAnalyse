package main

import (
	"context"
	"net"
	"strconv"

	"LogAnalyse/app/service/user/rpc/config"
	"LogAnalyse/app/service/user/rpc/dao"
	"LogAnalyse/app/service/user/rpc/initialize"
	"LogAnalyse/app/shared/consts"
	user "LogAnalyse/app/shared/kitex_gen/user/userservice"
	"LogAnalyse/app/shared/log"
	"LogAnalyse/app/shared/middleware"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
)

func main() {
	// init
	IP, Port := initialize.InitFlag()
	r, info := initialize.InitNacos(Port)
	db := initialize.InitMysql()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	userDao := dao.NewUser(db)
	impl := &UserServiceImpl{
		JWT: middleware.NewJWT(config.GlobalServerConfig.JWTInfo.SigningKey),
		Dao: userDao,
	}

	// Create new server.
	srv := user.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)

	err := srv.Run()

	if err != nil {
		log.Zlogger.Fatal(err.Error())
	}
}
