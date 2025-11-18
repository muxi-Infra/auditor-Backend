package main

import (
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/api/response"
	sdk "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/v2/server/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	g := r.Group("/api/v1")

	var mid = func(next sdk.HandlerFunc) sdk.HandlerFunc {
		return func(c *sdk.Context) (any, error) {
			fmt.Println("这是执行具体逻辑前的中间件")
			return next(c)
		}
	}

	chain := sdk.NewChain(mid)
	s := sdk.NewGinRegistrar(g)
	var c Controller

	s.WebHook("/webhook", chain, c.WebHook)

	r.Run(":8081")
}

type Controller struct{}

func (c *Controller) WebHook(g *gin.Context, req *request.HookPayload) (response.Resp, error) {
	fmt.Printf("请求体内容：%#v\n", req)
	return response.Resp{}, nil
}
