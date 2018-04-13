package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ProxyModeFlag = "mode" // 启动模式
)

type ProxyConfig struct {
	RootDir  string `mapstructure:"home"`
	HttpAddr string `mapstructure:"http_addr"`
	// receiv data
	HttpPipePop *HttpPipePopConfig `mapstructure:"http_pipe_pop"` // 数据接收
	// proc data mode
	Mode string `mapstructure:"mode"` // 目前支持两种模式，一种是forward，另外一种是persistence
	// persistence
	Persistence *ProxyPersistenceConfig `mapstructure:"persistence"` // 持久化数据配置
	// forward
	Forward *ProxyForwardConfig `mapstructure:"forward"` // 转发配置
	Debug   bool                `mapstructure:"debug"`
}

type ProxyPersistenceConfig struct {
	Engine    string `mapstructure:"engine"`     // 支持local influxdb
	LocalPath string `mapstructure:"local_path"` // 先只写到一个文件中，后面再优化
}

type ProxyForwardConfig struct {
	Pipe         string              `mapstructure:"pipe"`           // 转发传输管道 支持kafka、http
	HttpPipePush *HttpPipePushConfig `mapstructure:"http_pipe_push"` // HTTP管道配置
}

func DefaultProxyConfig() *ProxyConfig {
	config := &ProxyConfig{
		HttpAddr:    ":8070",
		Mode:        DefaultProxyMode(),
		HttpPipePop: DefaultHttpPipePopConfig(),
		Persistence: DefaultProxyPersistenceConfig(),
		Forward:     DefaultProxyForwardConfig(),
	}
	return config
}

func DefaultProxyMode() string {
	return "persistence"
}

func DefaultProxyPersistenceConfig() *ProxyPersistenceConfig {
	return &ProxyPersistenceConfig{
		Engine:    "local",                                                // 默认本地存储
		LocalPath: filepath.Join(DefaultDir, dataDir, "persistence.data"), // 本地存储默认路径
	}
}

func DefaultProxyForwardConfig() *ProxyForwardConfig {
	return &ProxyForwardConfig{
		Pipe:         "http",
		HttpPipePush: DefaultHttpPipePushConfig(),
	}
}

func ParseProxyConfig() *ProxyConfig {
	config := new(ProxyConfig)
	path := filepath.Join(RootDir, defaultProxyConfigFilePath)
	if err := writeConfig(path, DefaultProxyConfig()); err != nil {
		Exit("write default proxy config err:" + err.Error())
	}

	if err := decodeConfig(path, config); err != nil {
		Exit("parse proxy config err:" + err.Error())
	}

	config.HttpAddr = viper.GetString(HttpAddrFlag)
	config.Debug = viper.GetBool(DebugFlag)
	return config
}
