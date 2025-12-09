package config

import "github.com/muxi-Infra/auditor-Backend/pkg/viperx"

type MuxiAIType string

const (
	Ollama  MuxiAIType = "ollama"
	OpenAI  MuxiAIType = "open_ai"
	Alibaba MuxiAIType = "Alibaba"
)

type MuxiAI struct {
	Type            MuxiAIType `yaml:"type"`
	URL             string     `yaml:"url"`
	Model           string     `yaml:"model"`
	ApiKey          string     `yaml:"api_key"`
	AccessKeyID     string     `mapstructure:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string     `mapstructure:"access_key_secret" yaml:"access_key_secret"`
	Region          string     `yaml:"region"`
	Endpoint        string     `yaml:"endpoint"`
}

func NewMuxiAIConf(s *viperx.VipperSetting) *MuxiAI {
	var aiConf = &MuxiAI{}
	err := s.ReadSection("AI", &aiConf)
	if err != nil {
		panic(err)
	}
	return aiConf
}
