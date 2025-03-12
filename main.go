package main

import (
	"github.com/gin-gonic/gin"
	conf "muxi_auditor/config"
)

func main() {

	//TODO,改为从环境变量读取
	app := InitWebServer("./config/config-example.yaml")
	app.Run()

}

type App struct {
	r *gin.Engine
	c *conf.AppConf
}

func NewApp(r *gin.Engine, c *conf.AppConf) *App {
	return &App{
		r: r,
		c: c,
	}
}

// 启动
func (a *App) Run() {
	err := a.r.Run(a.c.Addr)
	if err != nil {
		panic(err)
	}
}
