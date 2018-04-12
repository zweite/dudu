package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type AgentConfig struct {
	RootDir       string           `mapstructure:"home"`
	Addr          string           `mapstructure:"addr"`          // 服务地址
	Compactor     string           `mapstructure:"compactor"`     // 压缩算法
	Pipe          string           `mapstructure:"pipe"`          // 传输管道
	HttpPipe      *HttpPipeConfig  `mapstructure:"httpPipe"`      // HTTP 管道配置
	BatchDuration int64            `mapstructure:"batchDuration"` // 批量传输最大间隔
	BatchLength   int              `mapstructure:"batchLength"`   // 批量传输最大长度
	Collects      []*CollectConfig `mapstructure:"collects"`      // 采集器配置
	Debug         bool             `mapstructure:"debug"`
}

type HttpPipeConfig struct {
	Addr string `mapstructure:"addr"`
	Auth string `mapstructure:"auth"`
}

type CollectConfig struct {
	Name     string // 采集器名称
	Duration int64  // 采集间隔
}

func DefaultAgentConfig() *AgentConfig {
	config := &AgentConfig{
		Addr:          ":8071",
		BatchDuration: 5,
		BatchLength:   100,
		Compactor:     "snappy",
		Pipe:          "http",
		HttpPipe: &HttpPipeConfig{
			Addr: "http://127.0.0.1:8070/collect",
			Auth: "midea",
		},
		Collects: []*CollectConfig{
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
		},
	}

	return config
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
	fmt.Printf("%+v\n", config)
	config.Addr = viper.GetString(AddrFlag)
	config.Debug = viper.GetBool(DebugFlag)
	return config
}
