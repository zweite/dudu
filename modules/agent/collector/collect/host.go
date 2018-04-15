package collect

import (
	"dudu/models"
	"dudu/modules/agent/collector"

	"github.com/shirou/gopsutil/host"
)

// 系统相关
type BootTime struct{}

// collect info
func (b *BootTime) Collect() (interface{}, error) {
	return host.BootTime()
}

// 采集数据类型
func (b *BootTime) Type() models.MetricType {
	return models.IndicatorMetricType
}

// collector name
func (b *BootTime) Name() string {
	return "BootTime"
}

type KernelVersion struct{}

// collect info
func (k *KernelVersion) Collect() (interface{}, error) {
	return host.KernelVersion()
}

// 采集数据类型
func (k *KernelVersion) Type() models.MetricType {
	return models.InfoMetricType
}

// collector name
func (k *KernelVersion) Name() string {
	return "KernelVersion"
}

type PlatformInfo struct{}

func (p *PlatformInfo) Collect() (interface{}, error) {
	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		return nil, err
	}

	return &models.PlatformInfo{
		Platform: platform,
		Family:   family,
		Version:  version,
	}, nil
}

func (p *PlatformInfo) Type() models.MetricType {
	return models.InfoMetricType
}

func (p *PlatformInfo) Name() string {
	return "PlatformInfo"
}

type VirtualizationInfo struct{}

func (v *VirtualizationInfo) Collect() (interface{}, error) {
	system, role, err := host.Virtualization()
	if err != nil {
		return nil, err
	}

	return &models.VirtualizationInfo{
		System: system,
		Role:   role,
	}, nil
}

func (v *VirtualizationInfo) Type() models.MetricType {
	return models.InfoMetricType
}

func (v *VirtualizationInfo) Name() string {
	return "VirtualizationInfo"
}

type HostInfo struct{}

func (h *HostInfo) Collect() (interface{}, error) {
	return host.Info()
}

func (h *HostInfo) Type() models.MetricType {
	return models.InfoMetricType
}

func (h *HostInfo) Name() string {
	return "HostInfo"
}

type Users struct{}

func (u *Users) Collect() (interface{}, error) {
	return host.Users()
}

func (u *Users) Type() models.MetricType {
	return models.InfoMetricType
}

func (u *Users) Name() string {
	return "Users"
}

type Uptime struct{}

func (u *Uptime) Collect() (interface{}, error) {
	return host.Uptime()
}

func (u *Uptime) Type() models.MetricType {
	return models.IndicatorMetricType
}

func (u *Uptime) Name() string {
	return "Uptime"
}

func init() {
	collector.RegisterCollector(new(BootTime))
	collector.RegisterCollector(new(KernelVersion))
	collector.RegisterCollector(new(PlatformInfo))
	collector.RegisterCollector(new(VirtualizationInfo))
	collector.RegisterCollector(new(HostInfo))
	collector.RegisterCollector(new(Users))
	collector.RegisterCollector(new(Uptime))
}
