package collect

import (
	"github.com/shirou/gopsutil/disk"

	"dudu/models"
	"dudu/modules/agent/collector"
)

// 磁盘相关
type Partitions struct{}

func (p *Partitions) Collect() (interface{}, error) {
	return disk.Partitions(true)
}

func (p *Partitions) Type() models.MetricType {
	return models.InfoMetricType
}

func (p *Partitions) Name() string {
	return "Partitions"
}

// 监控特定路径
type Usage struct {
	path string
}

func NewUsage(path string) *Usage {
	return &Usage{
		path: path,
	}
}

func (u *Usage) Collect() (interface{}, error) {
	return disk.Usage(u.path)
}

func (u *Usage) Type() models.MetricType {
	return models.InfoMetricType
}

func (u *Usage) Name() string {
	return u.path + "@Usage"
}

func init() {
	collector.RegisterCollector(new(Partitions))
}
