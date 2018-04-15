package proxy

import (
	"dudu/commons/log"
	"dudu/commons/util"
	"encoding/json"
	"fmt"

	"dudu/commons/compactor"
	"dudu/config"
	"dudu/models"
)

type Persistence struct {
	trigger      util.AutoTrigger // flush 长度
	logger       log.Logger
	cfg          *config.ProxyPersistenceConfig
	resCfg       *config.ResourceConfig
	compactorSet map[string]compactor.Compactor
	persistor    Persistor
	parser       *Parser
}

func NewPersistence(
	cfg *config.ProxyPersistenceConfig,
	resCfg *config.ResourceConfig,
	logger log.Logger) (persistence *Persistence, err error) {

	var persistor Persistor
	switch cfg.Engine {
	case "":
		fallthrough
	case "local":
		persistor, err = NewFilePersistor(cfg.LocalPath)
	default:
		return nil, fmt.Errorf("persistor not found 【%s】", cfg.Engine)
	}

	return &Persistence{
		cfg:          cfg,
		resCfg:       resCfg,
		logger:       logger,
		parser:       NewParser(logger),
		persistor:    persistor,
		compactorSet: compactor.GetCompactorSet(),
	}, err
}

func (p *Persistence) Stop() {
	if err := p.persistor.Flush(); err != nil {
		p.logger.Warnf("persistor flush err:%s", err.Error())
	}
}

// 存储需解压数据
func (p *Persistence) Proc(data []byte) (err error) {
	metric := new(models.MetricValue)
	if err = json.Unmarshal(data, metric); err != nil {
		return
	}

	if metric.Compactor != "" {
		compactor, ok := p.compactorSet[metric.Compactor]
		if !ok {
			return fmt.Errorf("not found compactor 【%s】", metric.Compactor)
		}

		value, err := compactor.Decode(metric.Value)
		if err != nil {
			return err
		}

		metric.Value = value
	}

	// WAL
	// 持久化
	data, err = json.Marshal(metric)
	if err == nil {
		// 写入数据
		if _, err = p.persistor.Write(append(data, '\n')); err == nil {
			// 刷盘
			p.trigger.IncTrigger(100, func() {
				err = p.persistor.Flush()
			})
		}
	}

	if err != nil {
		return
	}

	// 解析数据
	collectResults, err := p.parser.Parser(metric)
	if err != nil {
		return
	}

	for _, collectResult := range collectResults {
		if collectResult.Err != "" {
			p.logger.Warnf("endPoint:%s Metric:%s collect err:%s",
				metric.Endpoint, collectResult.Metric, collectResult.Err)
			continue
		}

		if err := p.procResult(metric, collectResult); err != nil {
			p.logger.Warnf("endPoint:%s Metric:%s proc collect result err:%s",
				metric.Endpoint, collectResult.Metric, err.Error())
		}
	}
	return
}

func (p *Persistence) procResult(metric *models.MetricValue, collectResult *models.CollectResult) error {
	switch collectResult.Type {
	case models.InfoMetricType:
		return p.procInfoResult(metric, collectResult)
	case models.IndicatorMetricType:
		return p.procIndicatorResult(metric, collectResult)
	}

	return fmt.Errorf("not found the metric type")
}

func (p *Persistence) procInfoResult(metric *models.MetricValue, collectResult *models.CollectResult) (err error) {
	p.logger.Infof("%s %+v", metric.Endpoint, collectResult)
	return
}

func (p *Persistence) procIndicatorResult(metric *models.MetricValue, collectResult *models.CollectResult) (err error) {
	p.logger.Infof("%s %+v", metric.Endpoint, collectResult)
	return
}
