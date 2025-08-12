package main

import (
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/webhook"
	"github.com/gin-gonic/gin"
)

func Handle(event string, data request.HookPayload) {
	fmt.Println("event:", event, "data:", data.Data)
}
func main() {
	r := gin.Default()
	l := webhook.NewListener(r, "0.0.0.0:8085", "/webhook", Handle)
	l.RegisterRoutes()
	l.Start()
}
