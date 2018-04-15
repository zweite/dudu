package collect

import (
	"encoding/json"

	"dudu/models"
	"dudu/modules/collector"

	"github.com/shirou/gopsutil/cpu"
)

// CPU相关
// CPU 核心统计
type CPUCount struct{}

func (c *CPUCount) Collect() (interface{}, error) {
	return cpu.Counts(true)
}

func (c *CPUCount) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (c *CPUCount) Unmarshal(data []byte) (interface{}, error) {
	var i int
	err := json.Unmarshal(data, &i)
	return i, err
}

func (c *CPUCount) Name() string {
	return "CPUCount"
}

func (c *CPUCount) Type() models.MetricType {
	return models.InfoMetricType
}

// CPU 信息
type CPUInfo struct{}

func (c *CPUInfo) Collect() (interface{}, error) {
	return cpu.Info()
}

func (c *CPUInfo) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (c *CPUInfo) Unmarshal(data []byte) (interface{}, error) {
	infoStats := make([]cpu.InfoStat, 0, 10)
	err := json.Unmarshal(data, &infoStats)
	return infoStats, err
}

func (c *CPUInfo) Name() string {
	return "CPUInfo"
}

func (c *CPUInfo) Type() models.MetricType {
	return models.InfoMetricType
}

// CPU 负载
type CPUTimes struct{}

func (c *CPUTimes) Collect() (interface{}, error) {
	return cpu.Times(true)
}

func (c *CPUTimes) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (c *CPUTimes) Unmarshal(data []byte) (interface{}, error) {
	timesStats := make([]cpu.TimesStat, 0, 10)
	err := json.Unmarshal(data, &timesStats)
	return timesStats, err
}

func (c *CPUTimes) Name() string {
	return "CPUTimes"
}

func (c *CPUTimes) Type() models.MetricType {
	return models.IndicatorMetricType
}

func init() {
	collector.RegisterCollector(new(CPUCount))
	collector.RegisterCollector(new(CPUInfo))
	collector.RegisterCollector(new(CPUTimes))

	collector.RegisterDefaultCollector(new(CPUCount))
	collector.RegisterDefaultCollector(new(CPUInfo))
	collector.RegisterDefaultCollector(new(CPUTimes))
}
