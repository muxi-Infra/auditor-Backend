package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/auditor-Backend/api/request"
	"github.com/muxi-Infra/auditor-Backend/sdk/v1/webhook"
)

func Handle(event string, data request.HookPayload) {
	fmt.Println("event:", event, "data:", data.Data)
}
func main() {
	r := gin.Default()
	l := webhook.NewListener(r, "0.0.0.0:8085", "/audit", Handle)
	l.RegisterRoutes()
	l.Start()
}
