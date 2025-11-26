package config

type Config struct {
	ApiKey         string
	ConnectTimeout int    // 单位毫秒，默认值为5000毫秒即5秒
	Region         string // 服务端所在地址，写到版本号即可
}
