package main

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/keyget"
	"github.com/gin-gonic/gin"
)

func main() {
	//ac:="cli_26JsjTlJcDdcmWKs"
	//sc:="26JsjTlJcDdcmWKsUz6mIrWHm9UTHYcb"
	e := gin.Default()
	keyget.DefaultServe(e, "localhost:8081", "/test").Run()

}
