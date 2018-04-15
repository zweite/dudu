package config

import "path/filepath"

type ResourceConfig struct {
	Redis []*DB `toml:"redis,omitempty" json:",omitempty"`
	Etcd  []*DB `toml:"etcd,omitempty" json:",omitempty"`
}

type DB struct {
	Host   string `toml:"host,omitempty" json:",omitempty"`
	Port   int    `toml:"port,omitempty" json:",omitempty"`
	User   string `toml:"user,omitempty" json:",omitempty"`
	Passwd string `toml:"passwd,omitempty" json:",omitempty"`
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
