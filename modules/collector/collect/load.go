package collect

import (
	"encoding/json"

	"dudu/models"
	"dudu/modules/collector"

	"github.com/shirou/gopsutil/load"
)

// 系统负载
type Avg struct{}

func (a *Avg) Collect() (interface{}, error) {
	return load.Avg()
}

func (a *Avg) Marshal(res interface{}) ([]byte, error) {
	return json.Marshal(res)
}

func (a *Avg) Unmarshal(data []byte) (interface{}, error) {
	avg := new(load.AvgStat)
	err := json.Unmarshal(data, avg)
	return avg, err
}

func (a *Avg) Type() models.MetricType {
	return models.IndicatorMetricType
}

func (a *Avg) Name() string {
	return "Avg"
}

func init() {
	collector.RegisterCollector(new(Avg))

	collector.RegisterDefaultCollector(new(Avg))
}
