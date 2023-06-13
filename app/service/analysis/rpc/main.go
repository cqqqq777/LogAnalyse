package main

import (
	"context"
	"net"
	"strconv"

	"LogAnalyse/app/service/analysis/rpc/config"
	"LogAnalyse/app/service/analysis/rpc/initialize"
	"LogAnalyse/app/service/analysis/rpc/pkg"
	"LogAnalyse/app/shared/consts"
	analysis "LogAnalyse/app/shared/kitex_gen/analyse/analyseservice"
	"LogAnalyse/app/shared/log"

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
	minioClient := initialize.InitMinio()
	producer, err := pkg.NewPublisher()
	if err != nil {
		log.Zlogger.Fatal("new nsq producer failed err:" + err.Error())
	}
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	impl := &AnalyseServiceImpl{
		Producer:     producer,
		MinioManager: pkg.NewMinioManager(minioClient),
	}

	srv := analysis.NewServer(impl,
		server.WithServiceAddr(utils.NewNetAddr(consts.TCP, net.JoinHostPort(IP, strconv.Itoa(Port)))),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithLimit(&limit.Option{MaxConnections: 2000, MaxQPS: 500}),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GlobalServerConfig.Name}),
	)

	err = srv.Run()

	if err != nil {
		log.Zlogger.Fatal(err.Error())
	}
}
