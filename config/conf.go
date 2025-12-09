package config

import (
	"github.com/muxi-Infra/auditor-Backend/pkg/viperx"
)

type AppConf struct {
	Addr string `yaml:"addr"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secretKey"` //秘钥
	Timeout   int    `yaml:"timeout"`   //过期时间
}

type DBConfig struct {
	Dsn string `yaml:"dsn"`
}

type CacheConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type OAuthConfig struct {
	Addr         string `yaml:"addr"`
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
}

type LogConfig struct {
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"maxSize"`    // 每个日志文件的最大大小，单位：MB
	MaxBackups int    `yaml:"maxBackups"` // 保留旧日志文件的最大个数
	MaxAge     int    `yaml:"maxAge"`     // 保留旧日志文件的最大天数
	Compress   int    `yaml:"compress"`   // 是否压缩旧的日志文件
}

type PrometheusConfig struct {
	Namespace string `yaml:"namespace"` //项目名称

	RouterCounter struct {
		Name string `yaml:"name"`
		Help string `yaml:"help"`
	} `yaml:"routerCounter"`

	ActiveConnections struct {
		Name string `yaml:"name"`
		Help string `yaml:"help"`
	} `yaml:"activeConnections"`

	DurationTime struct {
		Name string `yaml:"name"`
		Help string `yaml:"help"`
	} `yaml:"durationTime"`
}

type MiddlewareConf struct {
	AllowedOrigins []string `yaml:"allowedOrigins"`
}
type QiNiuYunConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
}

type KafkaConfig struct {
	Addr     []string `yaml:"addr"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
}

func NewAppConf(s *viperx.VipperSetting) *AppConf {
	var appConf = &AppConf{}
	err := s.ReadSection("app", appConf)
	if err != nil {
		return nil
	}
	return appConf
}

func NewJWTConf(s *viperx.VipperSetting) *JWTConfig {
	var jwtConf = &JWTConfig{}
	err := s.ReadSection("jwt", jwtConf)
	if err != nil {
		return nil
	}
	return jwtConf
}

func NewDBConf(s *viperx.VipperSetting) *DBConfig {
	var dbConf = &DBConfig{}
	err := s.ReadSection("db", dbConf)
	if err != nil {
		return nil
	}
	return dbConf
}

func NewOAuthConf(s *viperx.VipperSetting) *OAuthConfig {
	var oauthConf = &OAuthConfig{}
	err := s.ReadSection("oauth", oauthConf)
	if err != nil {
		return nil
	}
	return oauthConf
}

func NewCacheConf(s *viperx.VipperSetting) *CacheConfig {
	var cacheConf = &CacheConfig{}
	err := s.ReadSection("cache", cacheConf)
	if err != nil {
		return nil
	}
	return cacheConf
}

func NewLogConf(s *viperx.VipperSetting) *LogConfig {
	var logConf = &LogConfig{}
	err := s.ReadSection("log", logConf)
	if err != nil {
		return nil
	}
	return logConf
}

func NewPrometheusConf(s *viperx.VipperSetting) *PrometheusConfig {
	var prometheusConf = &PrometheusConfig{}
	err := s.ReadSection("prometheus", prometheusConf)
	if err != nil {
		return nil
	}
	return prometheusConf
}

func NewMiddleWareConf(s *viperx.VipperSetting) *MiddlewareConf {
	var middlewareConf = &MiddlewareConf{}
	err := s.ReadSection("middleware", middlewareConf)
	if err != nil {
		return nil
	}
	return middlewareConf
}
func NewQiniuConf(s *viperx.VipperSetting) *QiNiuYunConfig {
	var qiniuConf = &QiNiuYunConfig{}
	err := s.ReadSection("QiNiuYun", qiniuConf)
	if err != nil {
		return nil
	}
	return qiniuConf
}

func NewKafkaConf(s *viperx.VipperSetting) *KafkaConfig {
	var kafkaConf = &KafkaConfig{}

	err := s.ReadSection("kafka", kafkaConf)
	if err != nil {
		return nil
	}

	return kafkaConf
}
