package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ProxyModeFlag = "mode" // 启动模式
)

type ProxyConfig struct {
	RootDir  string `toml:"home,omitempty"`
	HttpAddr string `toml:"http_addr,omitempty"`
	// receiv data
	HttpPipePop *HttpPipePopConfig `toml:"http_pipe_pop,omitempty"` // 数据接收
	// proc data mode
	Mode string `toml:"mode,omitempty"` // 目前支持两种模式，一种是forward，另外一种是persistence
	// persistence
	Persistence *ProxyPersistenceConfig `toml:"persistence,omitempty"` // 持久化数据配置
	// forward
	Forward *ProxyForwardConfig `toml:"forward,omitempty"` // 转发配置
	Debug   bool                `toml:"debug,omitempty"`
}

type ProxyPersistenceConfig struct {
	WALEngine string `toml:"wal_engine,omitempty"` // 支持local
	LocalPath string `toml:"local_path,omitempty"` // 先只写到一个文件中，后面再优化

	InfoStorage      string `toml:"info_storage,omitempty`       // 信息存储 支持mongo
	IndicatorStorage string `toml:"indicator_storage,omitempty"` // 指标存储 支持
}

type ProxyForwardConfig struct {
	Pipe         string              `toml:"pipe,omitempty"`           // 转发传输管道 支持kafka、http
	HttpPipePush *HttpPipePushConfig `toml:"http_pipe_push,omitempty"` // HTTP管道配置
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
		WALEngine:        "local",                                                // 默认本地存储
		LocalPath:        filepath.Join(DefaultDir, dataDir, "persistence.data"), // 本地存储默认路径
		InfoStorage:      "mongo",                                                // mongodb
		IndicatorStorage: "influx",                                               // influxdb
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
