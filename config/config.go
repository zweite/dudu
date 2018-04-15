package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"dudu/commons/cli"
	"dudu/commons/util"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 配置

var (
	DefaultDir = os.ExpandEnv("${GOPATH}/src/dudu/config")
	RootDir    = DefaultDir
)

const (
	configDir          = "config"
	dataDir            = "data"
	baseConfigName     = "config.toml"
	proxyConfigName    = "proxy.toml"
	agentConfigName    = "agent.toml"
	resourceConfigName = "resource.toml"

	HttpAddrFlag = "http_addr"
	DebugFlag    = "debug"
)

var (
	defaultBaseConfigFilePath     = filepath.Join(configDir, baseConfigName)
	defaultProxyConfigFilePath    = filepath.Join(configDir, proxyConfigName)
	defaultAgentConfigFilePath    = filepath.Join(configDir, agentConfigName)
	defaultResourceConfigFilePath = filepath.Join(configDir, resourceConfigName)
)

type Config struct {
	BaseConfig
	Proxy    *ProxyConfig    `toml:"proxy,omitempty" json:",omitempty"`
	Agent    *AgentConfig    `toml:"agent,omitempty" json:",omitempty"`
	Resource *ResourceConfig `toml:"resource,omitempty" json:",omitempty"`
}

func (c *Config) SetRoot(root string) {
	c.Proxy.RootDir = root
}

type BaseConfig struct {
	RootDir  string `toml:"home,omitempty"`
	HostName string `toml:"-"`                // 节点主机名
	IPSeg    []byte `toml:"ip_seg,omitempty"` // 本节点所属IP段
	IP       string `toml:"-"`                // 本节点IP, 自动填充
	LogLevel string `toml:"log_level,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Agent:      DefaultAgentConfig(),
		Proxy:      DefaultProxyConfig(),
		Resource:   DefaultResourceConfig(),
	}
}

func ParseConfig() *Config {
	RootDir = viper.GetString(cli.ConfigFlag)
	return &Config{
		BaseConfig: ParseBaseConfig(),
		Agent:      ParseAgentConfig(),
		Proxy:      ParseProxyConfig(),
		Resource:   ParseResourceConfig(),
	}
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		IPSeg:    DefaultIPSeg(),
		HostName: getHostName(),
		LogLevel: DefaultLogLevel(),
	}
}

func ParseBaseConfig() BaseConfig {
	var config BaseConfig
	path := filepath.Join(RootDir, defaultBaseConfigFilePath)
	if err := writeConfig(path, DefaultBaseConfig()); err != nil {
		Exit("write default base config err:" + err.Error())
	}

	if err := decodeConfig(path, &config); err != nil {
		Exit("parse base config err:" + err.Error())
	}

	config.RootDir = RootDir
	config.IP = getLocalIP(config.IPSeg)
	config.HostName = getHostName()
	return config
}

func getHostName() string {
	hostname, _ := os.Hostname()
	return hostname
}

// 获取本地IP
func getLocalIP(ipSeg []byte) string {
	util.SetLocalIPSegment(ipSeg...)
	ip, err := util.LocalIP()
	if err != nil {
		Exit("get local ip err:" + err.Error())
	}
	return ip.String()
}

func DefaultIPSeg() []byte {
	return []byte{10, 73}
}

func DefaultLogLevel() string {
	return "info"
}

func DefaultLogLevelInt() logrus.Level {
	level, err := logrus.ParseLevel(DefaultLogLevel())
	if err != nil {
		Exit(err.Error())
	}
	return level
}

// 写入配置文件
func writeConfig(path string, defaultConfig interface{}) (err error) {
	if !util.FileExists(path) {
		dir := filepath.Dir(path)
		if err = util.EnsureDir(dir, 0755); err != nil {
			return
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		defer file.Close()
		encoder := toml.NewEncoder(file)
		return encoder.Encode(defaultConfig)
	}
	return nil
}

// 解码配置文件
func decodeConfig(path string, config interface{}) (err error) {
	data, err := readFile(path)
	if err != nil {
		return
	}

	if _, err = toml.Decode(string(data), config); err != nil {
		return
	}

	return
}

// 读取文件(有可能从配置中心读取)
func readFile(filePath string) (data []byte, err error) {
	data, err = ioutil.ReadFile(filePath)
	return
}

func Exit(s string) {
	fmt.Printf(s + "\n")
	os.Exit(1)
}
