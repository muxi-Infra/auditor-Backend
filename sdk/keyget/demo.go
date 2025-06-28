package keyget

import "github.com/gin-gonic/gin"

func main() {
	//ac:="cli_26JsjTlJcDdcmWKs"
	//sc:="26JsjTlJcDdcmWKsUz6mIrWHm9UTHYcb"
	e := gin.Default()
	DefaultServe(e, "localhost:8081", "/test").Run()

}
