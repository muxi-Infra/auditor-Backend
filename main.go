package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	conf "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
)

func main() {
	// 可选加载 .env（不存在则忽略）；容器内建议用环境变量或 env_file 注入
	_ = godotenv.Load()
	app := InitWebServer(os.Getenv("CONFIG_PATH"))
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
