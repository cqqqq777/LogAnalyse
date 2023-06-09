package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type NsqConfig struct {
	Host          string `mapstructure:"host" json:"host"`
	Port          int    `mapstructure:"port" json:"port"`
	ConsumerTopic string `mapstructure:"consumer_topic" json:"consumer_topic"`
	Channel       string `mapstructure:"channel" json:"channel"`
}

type ServerConfig struct {
	Name        string       `mapstructure:"name" json:"name"`
	Host        string       `mapstructure:"host" json:"host"`
	Port        int          `mapstructure:"port" json:"port"`
	JWTInfo     JWTConfig    `mapstructure:"jwt" json:"jwt"`
	OtelInfo    OtelConfig   `mapstructure:"otel" json:"otel"`
	NsqInfo     NsqConfig    `mapstructure:"nsq" json:"nsq"`
	RedisInfo   RedisConfig  `mapstructure:"redis" json:"redis"`
	UserSrvInfo RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	FileSrvInfo RPCSrvConfig `mapstructure:"file_srv" json:"file_srv"`
	JobSrvInfo  RPCSrvConfig `mapstructure:"job_srv" json:"job_srv"`
}

type RPCSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
