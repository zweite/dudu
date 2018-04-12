package collect

import (
	"dudu/modules/agent/collector"

	"github.com/shirou/gopsutil/cpu"
)

// CPU相关
// CPU 核心统计
type CPUCount struct{}

func (c *CPUCount) Collect() (interface{}, error) {
	return cpu.Counts(true)
}

func (c *CPUCount) Name() string {
	return "CPUCount"
}

// CPU 信息
type CPUInfo struct{}

func (c *CPUInfo) Collect() (interface{}, error) {
	return cpu.Info()
}

func (c *CPUInfo) Name() string {
	return "CPUInfo"
}

// CPU 负载
type CPUTimes struct{}

func (c *CPUTimes) Collect() (interface{}, error) {
	return cpu.Times(true)
}

func (c *CPUTimes) Name() string {
	return "CPUTimes"
}

func init() {
	collector.RegisterCollector(new(CPUCount))
	collector.RegisterCollector(new(CPUInfo))
	collector.RegisterCollector(new(CPUTimes))
}
