package models

import (
	"encoding/json"
)

type MetricType int

const (
	InfoMetricType MetricType = iota
	IndicatorMetricType
)

// 采集结果
type CollectResult struct {
	Metric   string          `json:",omitempty"` // 采集器名称
	Value    json.RawMessage `json:",omitempty"` // 采集信息
	RelValue interface{}     `json:",omitempty"` // 解析后的值（对象或者数组）
	Type     MetricType      `json:",omitempty"` // 0 信息类型 1 指标类型
	Version  int64           `json:",omitempty"` // 采集数据版本
	Err      string          `json:",omitempty"` // 采集出错信息
}
