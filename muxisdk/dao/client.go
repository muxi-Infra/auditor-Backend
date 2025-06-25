package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type OrmClient struct {
	DB *gorm.DB
}

func NewOrmClient(dsn string) *OrmClient {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return &OrmClient{
		DB: db,
	}
}
func (o *OrmClient) InitTable(fields ...any) {
	err := o.DB.AutoMigrate(fields...)
	if err != nil {
		panic(err)
	}
}
