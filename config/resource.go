package config

import "path/filepath"

type ResourceConfig struct {
	Redis  map[string][]*DBConfig `toml:"redis,omitempty" json:",omitempty"`  // redis
	Etcd   map[string][]*DBConfig `toml:"etcd,omitempty" json:",omitempty"`   // etcd
	Mongo  map[string][]*DBConfig `toml:"mongo,omitempty" json:",omitempty"`  // mongo
	Influx *DBConfig              `toml:"influx,omitempty" json:",omitempty"` // influx 暂不支持集群方案
}

type DBConfig struct {
	Host    string `toml:"host,omitempty" json:",omitempty"`
	Port    int    `toml:"port,omitempty" json:",omitempty"`
	Db      string `toml:"db,omitempty" json:",omitempty"`
	User    string `toml:"user,omitempty" json:",omitempty"`
	Passwd  string `toml:"passwd,omitempty" json:",omitempty"`
	Options string `toml:"options,omitempty" json:",omitempty"`
}

func DefaultResourceConfig() *ResourceConfig {
	config := &ResourceConfig{
		Redis: map[string][]*DBConfig{
			"base": []*DBConfig{
				&DBConfig{
					Host: "127.0.0.1",
					Port: 6379,
				},
			},
		},
		Etcd: map[string][]*DBConfig{
			"base": []*DBConfig{
				&DBConfig{
					Host: "127.0.0.1",
					Port: 4010,
				},
				&DBConfig{
					Host: "127.0.0.1",
					Port: 4011,
				},
				&DBConfig{
					Host: "127.0.0.1",
					Port: 4012,
				},
			},
		},
		Mongo: map[string][]*DBConfig{
			"base": []*DBConfig{
				&DBConfig{
					Host: "127.0.0.1",
					Port: 27017,
					Db:   "dudu",
				},
			},
		},
		Influx: &DBConfig{
			Host: "127.0.0.1",
			Port: 8086,
			Db:   "dudu",
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
