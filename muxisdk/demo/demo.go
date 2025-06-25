package main

import (
	"context"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/muxisdk/dao"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `gorm:"column:name;unique;not null"` // 指定列名为 name，唯一约束且不能为空
	Email  string `gorm:"column:email;not null"`
	Avatar int    `gorm:"column:avatar"`
}

func main() {
	//创建Client并初始化表
	c := dao.NewOrmClient("root:chenhaoqi318912@tcp(60.205.12.92:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local")
	c.InitTable(&User{})
	var co = context.Background()
	var user = User{Name: "c", Email: "chenhaoqi318912@qq.com", Avatar: 12}
	//增添数据，需要增添什么就传什么
	err := c.Create(co, &user)
	if err != nil {
		fmt.Println(err)
	}
	//获取数据，需要什么数据就传什么数据类型
	var user1 User
	c.Get(co, &user1, dao.Key{"Avatar", 12})
	err = c.Update(co, &User{
		Name:   "",
		Email:  "",
		Avatar: 122,
	}, dao.Key{"Name", "c"})
	if err != nil {
		fmt.Println(err)
	}
	c.Get(co, &user1, dao.Key{"Avatar", 122})
	fmt.Println(user1)
	err = c.Delete(co, &user1, dao.Key{"Avatar", 122})
	if err != nil {
		fmt.Println(err)
	}
}
