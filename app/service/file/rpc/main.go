package main

import (
	"context"
	"net"
	"strconv"

	"LogAnalyse/app/service/file/rpc/config"
	"LogAnalyse/app/service/file/rpc/initialize"
	"LogAnalyse/app/service/file/rpc/pkg"
	"LogAnalyse/app/shared/consts"
	file "LogAnalyse/app/shared/kitex_gen/file/fileservice"
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
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.GlobalServerConfig.Name),
		provider.WithExportEndpoint(config.GlobalServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	minioManager := pkg.NewMinioManager(minioClient)
	impl := &FileServiceImpl{Minio: minioManager}

	// Create new server.
	srv := file.NewServer(impl,
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
