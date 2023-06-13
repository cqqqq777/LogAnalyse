package initialize

import (
	"LogAnalyse/app/service/analysis/rpc/config"
	"LogAnalyse/app/shared/consts"
	"LogAnalyse/app/shared/log"
	"LogAnalyse/app/shared/tools"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/utils"
	nacos "github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"net"
	"strconv"
)

func InitNacos(Port int) (registry.Registry, *registry.Info) {
	v := viper.New()
	v.SetConfigFile(consts.AnalysisConfigFile)
	if err := v.ReadInConfig(); err != nil {
		msg := fmt.Sprintf("viper read analysis nacos config failed err:%v", err)
		log.Zlogger.Fatal(msg)
	}
	if err := v.Unmarshal(&config.GlobalNacosConfig); err != nil {
		msg := fmt.Sprintf("viper unmarshal analysis nacos config failed err:%v", err)
		log.Zlogger.Fatal(msg)
	}

	//read config from nacos
	sc := []constant.ServerConfig{
		{
			IpAddr: config.GlobalNacosConfig.Host,
			Port:   config.GlobalNacosConfig.Port,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         config.GlobalNacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              consts.NacosLogDir,
		CacheDir:            consts.NacosCacheDir,
		LogLevel:            consts.NacosLogLevel,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		msg := fmt.Sprintf("create config client failed err:%v", err)
		log.Zlogger.Fatal(msg)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: config.GlobalNacosConfig.DataId,
		Group:  config.GlobalNacosConfig.Group,
	})
	if err != nil {
		msg := fmt.Sprintf("read config client failed err:%v", err)
		log.Zlogger.Fatal(msg)
	}
	err = sonic.Unmarshal([]byte(content), &config.GlobalServerConfig)
	if err != nil {
		msg := fmt.Sprintf("nacos config failed err:%v", err)
		log.Zlogger.Fatal(msg)
	}
	if config.GlobalServerConfig.Host == "" {
		config.GlobalServerConfig.Host, err = tools.GetLocalIPv4Address()
		if err != nil {
			msg := fmt.Sprintf("get local ipv4 addr failed err:%v", err)
			log.Zlogger.Fatal(msg)
		}
	}
	registryClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		msg := fmt.Sprintf("create registry client failed err:%v", err)
		log.Zlogger.Fatal(msg)
	}

	r := nacos.NewNacosRegistry(registryClient, nacos.WithGroup(config.GlobalNacosConfig.Group))
	sf, err := snowflake.NewNode(consts.NacosSnowflakeNode)
	if err != nil {
		msg := fmt.Sprintf("generate service name failed: %s", err)
		log.Zlogger.Fatal(msg)
	}
	info := &registry.Info{
		ServiceName: config.GlobalServerConfig.Name,
		Addr:        utils.NewNetAddr(consts.TCP, net.JoinHostPort(config.GlobalServerConfig.Host, strconv.Itoa(Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
	}
	return r, info
}
