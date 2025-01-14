package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"muxi_auditor/config"
	"muxi_auditor/pkg/logger"
	"muxi_auditor/repository/dao"
)

func InitDB(conf *config.DBConfig, l logger.Logger) *gorm.DB {

	db, err := gorm.Open(mysql.Open(conf.Addr), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			SlowThreshold: 0,
			LogLevel:      glogger.Info, // 以Debug模式打印所有Info级别能产生的gorm日志
		}),
	})
	if err != nil {
		panic(err)
	}
	//初始化所有表
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

// TODO 修改日志系统
func (g gormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{Key: "args", Val: i})
}
