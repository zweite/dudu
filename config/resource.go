package config

import "path/filepath"

type ResourceConfig struct {
	Redis []*DB `mapstructure:"redis" json:",omitempty"`
	Etcd  []*DB `mapstructure:"etcd" json:",omitempty"`
}

type DB struct {
	Host   string `mapstructure:"host" json:",omitempty"`
	Port   int    `mapstructure:"port" json:",omitempty"`
	User   string `mapstructure:"user" json:",omitempty"`
	Passwd string `mapstructure:"passwd" json:",omitempty"`
}

func DefaultResourceConfig() *ResourceConfig {
	config := &ResourceConfig{
		Redis: []*DB{
			&DB{
				Host: "127.0.0.1",
				Port: 6379,
			},
		},
		Etcd: []*DB{
			&DB{
				Host: "127.0.0.1",
				Port: 4010,
			},
			&DB{
				Host: "127.0.0.1",
				Port: 4011,
			},
			&DB{
				Host: "127.0.0.1",
				Port: 4012,
			},
		},
	}

	return config
}

func ParseResourceConfig() *ResourceConfig {
	config := new(ResourceConfig)
	path := filepath.Join(RootDir, defaultResourceConfigFilePath)
	if err := writeConfig(path, DefaultResourceConfig()); err != nil {
		Exit("write default resource config err:" + err.Error())
	}

	if err := decodeConfig(path, config); err != nil {
		Exit("parse resource config err:" + err.Error())
	}
	return config
}
