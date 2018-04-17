package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type AgentConfig struct {
	RootDir       string              `toml:"home,omitempty"`
	Compactor     string              `toml:"compactor,omitempty"`      // 压缩算法
	Pipe          string              `toml:"pipe,omitempty"`           // 传输管道
	HttpPipePush  *HttpPipePushConfig `toml:"http_pipe_push,omitempty"` // HTTP 管道配置
	BatchDuration int64               `toml:"batch_duration,omitempty"` // 批量传输最大间隔
	BatchLength   int                 `toml:"batch_length,omitempty"`   // 批量传输最大长度
	Collects      []*CollectConfig    `toml:"collects,omitempty"`       // 采集器配置
	Debug         bool                `toml:"debug,omitempty"`
}

type HttpPipePushConfig struct {
	Addr string `toml:"addr,omitempty"`
	Auth string `toml:"auth,omitempty"`
}

type HttpPipePopConfig struct {
	Auth    string `toml:"auth,omitempty"`
	Pattern string `toml:"pattern,omitempty"`
}

type CollectConfig struct {
	Name     string // 采集器名称
	Duration int64  // 采集间隔
}

func DefaultAgentConfig() *AgentConfig {
	config := &AgentConfig{
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

	// config.HttpAddr = getHttpAddr(config.HttpAddr, DefaultAgentConfig().HttpAddr)
	config.Debug = viper.GetBool(DebugFlag)
	return config
}
