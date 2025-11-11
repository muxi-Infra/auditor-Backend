package dao

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/ioc"
)

func Init() ItemDaoInterface {
	var conf1 = config.DBConfig{Dsn: "root:chenhaoqi318912@tcp(60.205.12.92:3306)/muxiAuditor?charset=utf8mb4&parseTime=True&loc=Local"}
	var conf2 = config.LogConfig{
		Path:       "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   1,
	}
	log := ioc.InitLogger(&conf2)
	db := ioc.InitDB(&conf1, log)
	return NewItemDao(db)
}
