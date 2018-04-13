package proxy

import (
	"dudu/commons/log"
	"encoding/json"
	"fmt"

	"dudu/commons/compactor"
	"dudu/config"
	"dudu/models"
)

type Persistence struct {
	logger       log.Logger
	cfg          *config.ProxyPersistenceConfig
	resCfg       *config.ResourceConfig
	compactorSet map[string]compactor.Compactor
	persistor    Persistor
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
		persistor:    persistor,
		compactorSet: compactor.GetCompactorSet(),
	}, err
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

	data, err = json.Marshal(metric)
	if err != nil {
		return
	}

	if _, err = p.persistor.Write(append(data, '\n')); err != nil {
		return
	}

	// 刷盘策略有问题，可能会导致磁盘碎片产生
	return p.persistor.Flush()
}
