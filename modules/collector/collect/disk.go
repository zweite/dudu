package collect

import (
	"encoding/json"

	"dudu/models"
	"dudu/modules/collector"

	"github.com/shirou/gopsutil/disk"
)

// 磁盘相关
type Partitions struct{}

func (p *Partitions) Collect() (interface{}, error) {
	return disk.Partitions(true)
}

func (p *Partitions) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (p *Partitions) Unmarshal(data []byte) (interface{}, error) {
	partitionStats := make([]disk.PartitionStat, 0, 10)
	err := json.Unmarshal(data, &partitionStats)
	return partitionStats, err
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

func (u *Usage) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (u *Usage) Unmarshal(data []byte) (interface{}, error) {
	usageStat := new(disk.UsageStat)
	err := json.Unmarshal(data, usageStat)
	return usageStat, err
}

func (u *Usage) Type() models.MetricType {
	return models.InfoMetricType
}

func (u *Usage) Name() string {
	return u.path + "@Usage"
}

func init() {
	collector.RegisterCollector(new(Partitions))

	collector.RegisterDefaultCollector(new(Partitions))
	collector.RegisterDefaultCollector(NewUsage(""))
}
