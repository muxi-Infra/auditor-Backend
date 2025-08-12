package main

import (
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/sdk/keyget"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

const ApiKey = "remove.api_key"

type viperSetting struct {
	vp *viper.Viper
}

func newSetting(vp *viper.Viper) *viperSetting {
	return &viperSetting{vp: vp}
}
func (setting *viperSetting) SetApiKey(key string, value string) error {
	setting.vp.Set(key, value)
	if err := setting.vp.WriteConfig(); err != nil {
		log.Fatalf("写入配置错误: %v", err)
		return err
	}
	return nil
}
func main() {
	//ac:="cli_26JsjTlJcDdcmWKs"
	//sc:="26JsjTlJcDdcmWKsUz6mIrWHm9UTHYcb"
	e := gin.Default()
	vp := viper.New()
	vp.SetConfigFile("./demo.yaml") // 指定配置文件路径
	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}
	vp.WatchConfig()
	v := newSetting(vp)
	////default_sever
	//keyget.DefaultServe(e, "localhost:8081", "/test").Run()
	//write to your config.yaml sever
	keyget.ServeWriteToConfig(e, "localhost:8081", "/test", v, ApiKey).Run()
}
