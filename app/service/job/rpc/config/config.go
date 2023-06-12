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

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name            string          `mapstructure:"name" json:"name"`
	Host            string          `mapstructure:"host" json:"host"`
	OtelInfo        OtelConfig      `mapstructure:"otel" json:"otel"`
	MysqlInfo       MysqlConfig     `mapstructure:"mysql" json:"mysql"`
	AnalysisSrvInfo AnalysisService `mapstructure:"analysis_srv" json:"analysis_srv"`
	FileSrvInfo     FileService     `mapstructure:"file_srv" json:"file_srv"`
}

type AnalysisService struct {
	Name string `mapstructure:"name" json:"name"`
}

type FileService struct {
	Name string `mapstructure:"name" json:"name"`
}
