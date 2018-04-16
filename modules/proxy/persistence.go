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
	trigger            util.AutoTrigger // flush 长度
	logger             log.Logger
	cfg                *config.ProxyPersistenceConfig
	resCfg             *config.ResourceConfig
	compactorSet       map[string]compactor.Compactor
	walPersistor       Persistor     // 持久化日志
	infoPersistor      InfoPersistor // 信息持久化
	indicatorPersistor InfoPersistor // 指标持久化
	parser             *Parser
}

func NewPersistence(
	cfg *config.ProxyPersistenceConfig,
	resCfg *config.ResourceConfig,
	logger log.Logger) (persistence *Persistence, err error) {

	walPersistor, err := getWALPersistor(cfg)
	if err != nil {
		return
	}

	infoPersistor, err := getInfoPersistor(cfg.InfoStorage)
	if err != nil {
		return
	}

	indicatorPersistor, err := getInfoPersistor(cfg.IndicatorStorage)
	if err != nil {
		return
	}

	return &Persistence{
		cfg:                cfg,
		resCfg:             resCfg,
		logger:             logger,
		parser:             NewParser(logger),
		walPersistor:       walPersistor,
		infoPersistor:      infoPersistor,
		indicatorPersistor: indicatorPersistor,
		compactorSet:       compactor.GetCompactorSet(),
	}, err
}

func getInfoPersistor(engine string) (infoPersistor InfoPersistor, err error) {
	switch engine {
	case "":
		fallthrough
	case "mongo":
		infoPersistor, err = NewMongoPersistor("base")
	case "influx":
		resCfg := config.ParseResourceConfig()
		infoPersistor, err = NewInfluxPersistor(resCfg.Influx)
	default:
		return nil, fmt.Errorf("persistor not found 【%s】", engine)
	}
	return
}

func getWALPersistor(cfg *config.ProxyPersistenceConfig) (persistor Persistor, err error) {
	switch cfg.WALEngine {
	case "":
		fallthrough
	case "local":
		persistor, err = NewFilePersistor(cfg.LocalPath)
	default:
		return nil, fmt.Errorf("persistor not found 【%s】", cfg.WALEngine)
	}
	return
}

func (p *Persistence) Stop() {
	if err := p.walPersistor.Close(); err != nil {
		p.logger.Warnf("walPersistor close err:%s", err.Error())
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

	// 解析数据
	collectResults, err := p.parser.Parser(metric)
	if err != nil {
		return
	}

	for _, collectResult := range collectResults {
		// 写日志
		p.walPersistor.Write(
			metric.Endpoint, metric.HostName,
			collectResult.Metric, string(collectResult.Value),
			collectResult.Err, collectResult.Version)

		if collectResult.Err != "" {
			p.logger.Warnf("endPoint:%s Metric:%s collect err:%s",
				metric.Endpoint, collectResult.Metric, collectResult.Err)
			continue
		}

		if err := p.procResult(metric, collectResult); err != nil {
			p.logger.Warnf("endPoint:%s Metric:%s proc collect result data:%+v err:%s",
				metric.Endpoint, collectResult.Metric, collectResult, err.Error())
		}
	}

	// 刷盘
	p.trigger.IncTrigger(10, func() {
		err = p.walPersistor.Flush()
	})
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
	return p.infoPersistor.Save(
		metric.Endpoint, metric.HostName,
		collectResult.Metric, collectResult.RelValue, collectResult.Version,
	)
}

func (p *Persistence) procIndicatorResult(metric *models.MetricValue, collectResult *models.CollectResult) (err error) {
	return p.indicatorPersistor.Save(
		metric.Endpoint, metric.HostName,
		collectResult.Metric, collectResult.RelValue, collectResult.Version,
	)
}
