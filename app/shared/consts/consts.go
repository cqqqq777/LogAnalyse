package consts

// app context key constant
const (
	AccountID = "accountID"

	AccountIdentity = "accountIdentity"
)

// identity constant
const (
	UserIdentity = "user"
)

// expire time
const (
	TokenExpiredAt = 604800
)

// database constant
const (
	MysqlDns = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

// net constant
const (
	TCP             = "tcp"
	FreePortAddress = "localhost:0"
)

// ip and port constant
const (
	IPFlagName  = "ip"
	IPFlagValue = "0.0.0.0"
	IPFlagUsage = "address"

	PortFlagName  = "port"
	PortFlagUsage = "port"
)

// nacos constant
const (
	NacosLogDir   = "tmp/nacos/log"
	NacosCacheDir = "tmp/nacos/cache"
	NacosLogLevel = "debug"
)

// config file constant
const (
	ApiConfigFile = "./app/service/api/config.yaml"

	UserConfigFile     = "./app/service/user/rpc/config.yaml"
	FileConfigFile     = "./app/service/file/rpc/config.yaml"
	JobConfigFile      = "./app/service/job/rpc/config.yaml"
	AnalysisConfigFile = "./app/service/analysis/rpc/config.yaml"
)

const (
	NacosSnowflakeNode int64 = 0
)

const (
	ApiGroup = "API_GROUP"

	UserGroup    = "USER_GROUP"
	FileGroup    = "FILE_GROUP"
	JobGroup     = "JOB_GROUP"
	AnalyseGroup = "ANALYSIS_GROUP"
)
