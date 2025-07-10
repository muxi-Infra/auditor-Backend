package keyget

import (
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"

	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

//负责获取密钥

type KeyGet struct {
	Engine  *gin.Engine
	Addr    string
	Path    string
	Handler gin.HandlerFunc
}

func NewKeyGet(engine *gin.Engine, addr string, path string, handler gin.HandlerFunc) *KeyGet {
	return &KeyGet{
		Engine:  engine,
		Addr:    addr,
		Path:    path,
		Handler: handler,
	}
}
func (k *KeyGet) Serve() *KeyGet {
	k.Engine.POST(k.Path, k.Handler)
	return k
}
func (k *KeyGet) Run() error {
	return k.Engine.Run(k.Addr)
}

//提供的默认处理函数，写入到当前文件夹

func DefaultServe(engine *gin.Engine, addr string, path string) *KeyGet {
	handler := func(c *gin.Context) {
		var data request.ReturnApiKey
		err := c.ShouldBind(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filename := "./key.txt"
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("打开文件失败:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("api_key: %s\n", data.ApiKey))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}
	k := NewKeyGet(engine, addr, path, handler)
	k.Serve()
	return k
}
