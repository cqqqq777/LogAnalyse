package initialize

import (
	"LogAnalyse/app/service/api/config"
	"LogAnalyse/app/shared/consts"
	"LogAnalyse/app/shared/log"
	"LogAnalyse/app/shared/tools"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"net"
	"strconv"
)

func InitNacos() (registry.Registry, *registry.Info) {
	v := viper.New()
	v.SetConfigFile(consts.ApiConfigFile)
	if err := v.ReadInConfig(); err != nil {
		log.Zlogger.Fatal("read viper config failed: " + err.Error())
	}
	if err := v.Unmarshal(&config.GlobalNacosConfig); err != nil {
		log.Zlogger.Fatal("unmarshal err failed: " + err.Error())
	}
	log.Zlogger.Info(fmt.Sprintf("Config Info: %v", config.GlobalNacosConfig))

	// Read configuration information from nacos
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
		log.Zlogger.Fatal("create config client failed: " + err.Error())
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: config.GlobalNacosConfig.DataId,
		Group:  config.GlobalNacosConfig.Group,
	})
	if err != nil {
		log.Zlogger.Fatal("get config failed: " + err.Error())
	}

	err = sonic.Unmarshal([]byte(content), &config.GlobalServerConfig)
	if err != nil {
		log.Zlogger.Fatal("nacos config failed: %s" + err.Error())
	}

	if config.GlobalServerConfig.Host == "" {
		config.GlobalServerConfig.Host, err = tools.GetLocalIPv4Address()
		if err != nil {
			log.Zlogger.Fatal("get localIpv4Addr failed:%s" + err.Error())
		}
	}

	registryClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Zlogger.Fatal("create registryClient err: " + err.Error())
	}

	r := nacos.NewNacosRegistry(registryClient, nacos.WithRegistryGroup(consts.ApiGroup))

	sf, err := snowflake.NewNode(2)
	if err != nil {
		log.Zlogger.Fatal("generate service name failed: " + err.Error())
	}
	info := &registry.Info{
		ServiceName: config.GlobalServerConfig.Name,
		Addr: utils.NewNetAddr(consts.TCP, net.JoinHostPort(config.GlobalServerConfig.Host,
			strconv.Itoa(config.GlobalServerConfig.Port))),
		Tags: map[string]string{
			"ID": sf.Generate().Base36(),
		},
		Weight: registry.DefaultWeight,
	}

	return r, info
}
