package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type ProxyConfig struct {
	RootDir string `mapstructure:"home"`
	Addr    string `mapstructure:"addr"`
	Debug   bool   `mapstructure:"debug"`
}

func DefaultProxyConfig() *ProxyConfig {
	config := &ProxyConfig{
		Addr: ":8070",
	}
	return config
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

	config.Addr = viper.GetString(AddrFlag)
	config.Debug = viper.GetBool(DebugFlag)
	return config
}
