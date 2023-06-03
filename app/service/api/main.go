// Code generated by hertz generator.

package main

import (
	"fmt"

	"LogAnalyse/app/service/api/config"
	"LogAnalyse/app/service/api/initialize"
	"LogAnalyse/app/service/api/initialize/rpc"

	"github.com/cloudwego/hertz/pkg/app/server"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/pprof"
)

func main() {
	// init
	r, info := initialize.InitNacos()
	tracer, cfg := hertztracing.NewServerTracer()
	rpc.Init()

	// create a new server
	h := server.New(
		tracer,
		server.WithHostPorts(fmt.Sprintf(":%d", config.GlobalServerConfig.Port)),
		server.WithRegistry(r, info),
		server.WithHandleMethodNotAllowed(true),
	)

	// use pprof & tracer mw
	pprof.Register(h)
	h.Use(hertztracing.ServerMiddleware(cfg))
	register(h)
	h.Spin()
}
