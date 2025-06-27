package main

import (
	conf "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	godotenv.Load(".env")
	err := a.r.Run(a.c.Addr)

	if err != nil {
		panic(err)
	}
}
