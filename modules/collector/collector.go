package collector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"dudu/commons/log"
	"dudu/config"
	"dudu/models"
)

var (
	__Collectors           map[string]Collector
	__DefaultCollectorsSet map[string]Collector
)

func init() {
	__Collectors = make(map[string]Collector)
	__DefaultCollectorsSet = make(map[string]Collector)
}

// 采集器
type Collector interface {
	Collect() (interface{}, error)         // collect info
	Marshal(interface{}) ([]byte, error)   // 编码
	Unmarshal([]byte) (interface{}, error) // 解码
	Type() models.MetricType               // 采集数据类型
	Name() string                          // collector name
}

type CollectorManager struct {
	logger            log.Logger
	wg                *sync.WaitGroup
	mux               *sync.RWMutex
	collectsConfig    []*config.CollectConfig
	collectResultChan chan *models.CollectResult
	stopChans         map[string]chan struct{}
}

func NewCollectorManager(logger log.Logger, cfgs []*config.CollectConfig) *CollectorManager {
	return &CollectorManager{
		collectsConfig:    cfgs,
		logger:            logger,
		wg:                new(sync.WaitGroup),
		mux:               new(sync.RWMutex),
		stopChans:         make(map[string]chan struct{}),
		collectResultChan: make(chan *models.CollectResult, 100),
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
func (c *CollectorManager) Run() <-chan *models.CollectResult {
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
				var msg json.RawMessage
				var errMsg string
				if err != nil {
					errMsg = err.Error()
				} else {
					msg, err = collector.Marshal(result)
					if err != nil {
						errMsg = err.Error()
					}
				}
				c.collectResultChan <- &models.CollectResult{
					Metric:  collector.Name(),
					Value:   msg,
					Version: time.Now().UnixNano() / int64(time.Millisecond),
					Err:     errMsg,
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

// 注册全部类型采集器，以作为编码解码
func RegisterDefaultCollector(collectors ...Collector) {
	for _, collector := range collectors {
		if _, ok := __DefaultCollectorsSet[collector.Name()]; ok {
			// 重复注册
			continue
		}
		__DefaultCollectorsSet[collector.Name()] = collector
	}
}

func MarshalResult(name string, res interface{}) ([]byte, error) {
	name = getDefaultCollectorName(name)
	collector, ok := __DefaultCollectorsSet[name]
	if !ok {
		return nil, fmt.Errorf("not found %s collector", name)
	}

	return collector.Marshal(res)
}

func UnmarshalResult(name string, data []byte) (interface{}, error) {
	name = getDefaultCollectorName(name)
	collector, ok := __DefaultCollectorsSet[name]
	if !ok {
		return nil, fmt.Errorf("not found %s collector", name)
	}
	return collector.Unmarshal(data)
}

func getDefaultCollectorName(name string) string {
	fields := strings.Split(name, "@")
	if len(fields) == 2 {
		name = "@" + fields[1]
	}

	return name
}

// 获取全部采集器
func GetCollectors() []Collector {
	collectors := make([]Collector, 0, len(__Collectors))
	for _, collector := range __Collectors {
		collectors = append(collectors, collector)
	}
	return collectors
}
