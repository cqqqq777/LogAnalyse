package main

import (
	"LogAnalyse/app/shared/log"
	"context"
	"net"
	"strconv"

	"LogAnalyse/app/service/job/rpc/config"
	"LogAnalyse/app/service/job/rpc/dao"
	"LogAnalyse/app/service/job/rpc/initialize"
	"LogAnalyse/app/service/job/rpc/pkg"
	"LogAnalyse/app/shared/consts"
	job "LogAnalyse/app/shared/kitex_gen/job/jobservice"

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

	jobDao := dao.NewJob(db)
	fileManager := initialize.InitFile()
	analysisManager := initialize.InitAnalysis()

	impl := &JobServiceImpl{
		Dao:             jobDao,
		FileManager:     pkg.NewFileManager(fileManager),
		AnalysisManager: pkg.NewAnalysisManager(analysisManager),
	}

	// Create new server.
	srv := job.NewServer(impl,
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
