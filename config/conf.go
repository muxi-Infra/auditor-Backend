package config

import (
	"muxi_auditor/pkg/viperx"
)

type AppConf struct {
	Addr string `yaml:"addr"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secretKey"` //秘钥
	Timeout   int    `yaml:"timeout"`   //过期时间
}

type DBConfig struct {
	Addr string `yaml:"addr"`
}

type CacheConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type OAuthConfig struct {
	Addr         string `yaml:"addr"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type LogConfig struct {
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`    // 每个日志文件的最大大小，单位：MB
	MaxBackups int    `yaml:"max_backups"` // 保留旧日志文件的最大个数
	MaxAge     int    `yaml:"max_age"`     // 保留旧日志文件的最大天数
	Compress   int    `yaml:"compress"`    // 是否压缩旧的日志文件
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
