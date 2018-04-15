package collect

import (
	"encoding/json"

	"dudu/models"
	"dudu/modules/collector"

	"github.com/shirou/gopsutil/net"
)

// 网络相关
type Connections struct {
	kind string
}

func NewConnections(kind string) *Connections {
	if kind == "" {
		kind = "all"
	}
	return &Connections{
		kind: kind,
	}
}

// collect info
func (c *Connections) Collect() (interface{}, error) {
	return net.Connections(c.kind)
}

func (c *Connections) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (c *Connections) Unmarshal(data []byte) (interface{}, error) {
	connectionStats := make([]net.ConnectionStat, 0, 10)
	err := json.Unmarshal(data, &connectionStats)
	return connectionStats, err
}

// 采集数据类型
func (c *Connections) Type() models.MetricType {
	return models.InfoMetricType
}

// collector name
func (c *Connections) Name() string {
	return c.kind + "@Connections"
}

type Interfaces struct{}

// collect info
func (i *Interfaces) Collect() (interface{}, error) {
	return net.Interfaces()
}

func (i *Interfaces) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (i *Interfaces) Unmarshal(data []byte) (interface{}, error) {
	interfaceStats := make([]net.InterfaceStat, 0, 10)
	err := json.Unmarshal(data, &interfaceStats)
	return interfaceStats, err
}

// 采集数据类型
func (i *Interfaces) Type() models.MetricType {
	return models.InfoMetricType
}

// collector name
func (i *Interfaces) Name() string {
	return "Interfaces"
}

func init() {
	collector.RegisterCollector(new(Interfaces))

	collector.RegisterDefaultCollector(NewConnections(""))
	collector.RegisterDefaultCollector(new(Interfaces))
}
