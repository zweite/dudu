package collect

import (
	"dudu/models"
	"dudu/modules/agent/collector"

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
}
