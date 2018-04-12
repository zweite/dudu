package collector

import (
	"dudu/commons/log"
	"dudu/config"
	"sync"
	"time"
)

var (
	__Collectors map[string]Collector
)

func init() {
	__Collectors = make(map[string]Collector)
}

// 采集器
type Collector interface {
	Collect() (interface{}, error) // collect info
	Name() string                  // collector name
}

// 采集结果
type CollectResult struct {
	Metric  string      // 采集器名称
	Value   interface{} // 采集信息
	Version int64       // 采集数据版本
	Err     error       // 采集出错信息
}

type CollectorManager struct {
	logger            log.Logger
	wg                *sync.WaitGroup
	mux               *sync.RWMutex
	collectsConfig    []*config.CollectConfig
	collectResultChan chan *CollectResult
	stopChans         map[string]chan struct{}
}

func NewCollectorManager(logger log.Logger, cfgs []*config.CollectConfig) *CollectorManager {
	return &CollectorManager{
		collectsConfig:    cfgs,
		logger:            logger,
		wg:                new(sync.WaitGroup),
		mux:               new(sync.RWMutex),
		stopChans:         make(map[string]chan struct{}),
		collectResultChan: make(chan *CollectResult, 100),
	}
}

// 开启某一采集单元
func (c *CollectorManager) StartCollector(name string, duration time.Duration) bool {
	collector, ok := __Collectors[name]
	if !ok {
		return ok
	}

	return c.run(duration, collector)
}

// 开始采集
func (c *CollectorManager) Run() <-chan *CollectResult {
	validateCollectors := getValidateCollectors(c.collectsConfig)
	for _, collector := range validateCollectors {
		if cfg := c.getCollectorCfg(collector.Name()); cfg == nil {
			c.run(0, collector)
		} else {
			c.run(time.Second*time.Duration(cfg.Duration), collector)
		}
	}

	return c.collectResultChan
}

// 获取采集器配置
func (c *CollectorManager) getCollectorCfg(name string) *config.CollectConfig {
	for _, cfg := range c.collectsConfig {
		if cfg.Name == name {
			return cfg
		}
	}

	return nil
}

// 停止采集退出
func (c *CollectorManager) Stop() {
	c.mux.Lock()
	defer c.mux.Unlock()
	for name, stopChan := range c.stopChans {
		close(stopChan)
		delete(c.stopChans, name)
		c.logger.Infof("[%s] stoped", name)
	}
	c.wg.Wait()
	close(c.collectResultChan)
	return
}

// 停止采集某一单元
func (c *CollectorManager) StopCollector(name string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	stopChan, ok := c.stopChans[name]
	if !ok {
		return ok
	}

	delete(c.stopChans, name)
	close(stopChan)
	c.logger.Infof("[%s] stoped", name)
	return true
}

// 执行采集单元
func (c *CollectorManager) run(duration time.Duration, collector Collector) bool {
	c.mux.RLock()
	if _, ok := c.stopChans[collector.Name()]; ok {
		c.mux.RUnlock()
		return false
	}

	c.mux.RUnlock()
	c.logger.Infof("[%s] running", collector.Name())

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		if duration == 0 {
			duration = 30 * time.Second // 默认30s采集一次
		}

		stopChan := make(chan struct{})

		c.mux.Lock()
		c.stopChans[collector.Name()] = stopChan
		c.mux.Unlock()

		ticker := time.NewTicker(duration)
		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				result, err := collector.Collect()
				c.collectResultChan <- &CollectResult{
					Metric: collector.Name(),
					Value:  result,
					Err:    err,
				}
			}
		}
	}()
	return true
}

// 获取正在执行的单元
func (c *CollectorManager) GetRunningCollectors() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	names := make([]string, 0, len(c.stopChans))
	for name, _ := range c.stopChans {
		names = append(names, name)
	}
	return names
}

// 获取有效的采集器
func getValidateCollectors(cfgs []*config.CollectConfig) []Collector {
	validateCollectors := make([]Collector, 0, 10)
	for _, cfg := range cfgs {
		collector, ok := __Collectors[cfg.Name]
		if !ok {
			continue
		}
		validateCollectors = append(validateCollectors, collector)
	}
	return validateCollectors
}

// 注册采集器
func RegisterCollector(collectors ...Collector) {
	for _, collector := range collectors {
		if _, ok := __Collectors[collector.Name()]; ok {
			// 重复注册
			continue
		}
		__Collectors[collector.Name()] = collector
	}
}

// 获取全部采集器
func GetCollectors() []Collector {
	collectors := make([]Collector, 0, len(__Collectors))
	for _, collector := range __Collectors {
		collectors = append(collectors, collector)
	}
	return collectors
}