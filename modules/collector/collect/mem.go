package collect

import (
	"encoding/json"

	"dudu/models"
	"dudu/modules/collector"

	"github.com/shirou/gopsutil/mem"
)

// 内存相关
type SwapMemory struct{}

// collect info
func (s *SwapMemory) Collect() (interface{}, error) {
	return mem.SwapMemory()
}

func (s *SwapMemory) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (s *SwapMemory) Unmarshal(data []byte) (interface{}, error) {
	swapMemory := new(mem.SwapMemoryStat)
	err := json.Unmarshal(data, swapMemory)
	return swapMemory, err
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

func (v *VirtualMemory) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (v *VirtualMemory) Unmarshal(data []byte) (interface{}, error) {
	virtualMemoryStat := new(mem.VirtualMemoryStat)
	err := json.Unmarshal(data, virtualMemoryStat)
	return virtualMemoryStat, err
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

	collector.RegisterDefaultCollector(new(SwapMemory))
	collector.RegisterDefaultCollector(new(VirtualMemory))
}
