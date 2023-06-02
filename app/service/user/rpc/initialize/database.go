package initialize

import (
	"fmt"

	"LogAnalyse/app/service/user/rpc/config"
	"LogAnalyse/app/shared/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMysql init mysql
func InitMysql() *gorm.DB {
	dsn := config.GlobalServerConfig.GetMysqlDsn()
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		msg := fmt.Sprintf("init database failed err:%s", err.Error())
		log.Zlogger.Fatal(msg)
	}
	return db
}
