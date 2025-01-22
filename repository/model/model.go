package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `gorm:"column:name;unique;not null"` // 指定列名为 name，唯一约束且不能为空
	Email string `gorm:"column:email;not null"`
}
