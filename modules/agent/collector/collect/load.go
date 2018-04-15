package collect

import (
	"dudu/models"
	"dudu/modules/agent/collector"

	"github.com/shirou/gopsutil/load"
)

// 系统负载
type Avg struct{}

func (a *Avg) Collect() (interface{}, error) {
	return load.Avg()
}

func (a *Avg) Type() models.MetricType {
	return models.IndicatorMetricType
}

func (a *Avg) Name() string {
	return "Avg"
}

func init() {
	collector.RegisterCollector(new(Avg))
}
