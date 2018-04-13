package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type AgentConfig struct {
	RootDir       string              `mapstructure:"home"`
	HttpAddr      string              `mapstructure:"addr"`           // 服务地址
	Compactor     string              `mapstructure:"compactor"`      // 压缩算法
	Pipe          string              `mapstructure:"pipe"`           // 传输管道
	HttpPipePush  *HttpPipePushConfig `mapstructure:"http_pipe_push"` // HTTP 管道配置
	BatchDuration int64               `mapstructure:"batch_duration"` // 批量传输最大间隔
	BatchLength   int                 `mapstructure:"batch_length"`   // 批量传输最大长度
	Collects      []*CollectConfig    `mapstructure:"collects"`       // 采集器配置
	Debug         bool                `mapstructure:"debug"`
}

type HttpPipePushConfig struct {
	Addr string `mapstructure:"addr"`
	Auth string `mapstructure:"auth"`
}

type HttpPipePopConfig struct {
	Auth    string `mapstructure:"auth"`
	Pattern string `mapstructure:"pattern"`
}

type CollectConfig struct {
	Name     string // 采集器名称
	Duration int64  // 采集间隔
}

func DefaultAgentConfig() *AgentConfig {
	config := &AgentConfig{
		HttpAddr:      ":8071",
		BatchDuration: 5,
		BatchLength:   100,
		Compactor:     DefaultCompactor(),
		Pipe:          "http",
		HttpPipePush:  DefaultHttpPipePushConfig(),
		Collects:      DefaultCollectConfig(),
	}
	return config
}

func DefaultCompactor() string {
	return "gzip"
}

func DefaultHttpPipePopConfig() *HttpPipePopConfig {
	return &HttpPipePopConfig{
		Auth:    "midea",
		Pattern: "/collect",
	}
}

func DefaultHttpPipePushConfig() *HttpPipePushConfig {
	return &HttpPipePushConfig{
		Addr: "http://127.0.0.1:8070/collect",
		Auth: "midea",
	}
}

func DefaultCollectConfig() []*CollectConfig {
	return []*CollectConfig{
		&CollectConfig{
			Name:     "CPUCount",
			Duration: 5,
		},
		&CollectConfig{
			Name:     "CPUInfo",
			Duration: 5,
		},
		&CollectConfig{
			Name:     "CPUTimes",
			Duration: 5,
		},
	}
}

func ParseAgentConfig() *AgentConfig {
	config := new(AgentConfig)
	path := filepath.Join(RootDir, defaultAgentConfigFilePath)
	if err := writeConfig(path, DefaultAgentConfig()); err != nil {
		Exit("write default agent config err:" + err.Error())
	}

	if err := decodeConfig(path, config); err != nil {
		Exit("parse agent config err:" + err.Error())
	}
	config.HttpAddr = viper.GetString(HttpAddrFlag)
	config.Debug = viper.GetBool(DebugFlag)
	return config
}
