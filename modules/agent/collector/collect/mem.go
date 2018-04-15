package collect

import (
	"dudu/models"
	"dudu/modules/agent/collector"

	"github.com/shirou/gopsutil/mem"
)

// 内存相关
type SwapMemory struct{}

// collect info
func (s *SwapMemory) Collect() (interface{}, error) {
	return mem.SwapMemory()
}

// 采集数据类型
func (s *SwapMemory) Type() models.MetricType {
	return models.IndicatorMetricType
}

// collector name
func (s *SwapMemory) Name() string {
	return "SwapMemory"
}

type VirtualMemory struct{}

func (v *VirtualMemory) Collect() (interface{}, error) {
	return mem.VirtualMemory()
}

// 采集数据类型
func (v *VirtualMemory) Type() models.MetricType {
	return models.IndicatorMetricType
}

// collector name
func (v *VirtualMemory) Name() string {
	return "VirtualMemory"
}

func init() {
	collector.RegisterCollector(new(SwapMemory))
	collector.RegisterCollector(new(VirtualMemory))
}
