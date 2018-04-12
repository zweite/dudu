package agent

import (
	"dudu/commons/compactor"
	"encoding/json"
	"sync"

	"dudu/commons/log"
	"dudu/commons/util"
	"dudu/config"
	"dudu/modules/agent/collector"
)

type AgentNodeProvider func(*config.Config, log.Logger) (*AgentNode, error)

type AgentNode struct {
	wg           *sync.WaitGroup
	cfg          *config.Config
	logger       log.Logger
	collectorMag *collector.CollectorManager
	compactor    compactor.Compactor
	pipe         Pipe
	ctx          *AgentContext
}

func NewAgentNode(cfg *config.Config, logger log.Logger) (*AgentNode, error) {
	return &AgentNode{
		wg:     new(sync.WaitGroup),
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (app *AgentNode) Init() error {
	type f func() error

	initFuncs := []f{
		app.initPipe,      // 初始化管道
		app.initCompactor, // 初始化压缩器
		app.initCollect,   // 初始化采集器
	}

	for _, initFunc := range initFuncs {
		if err := initFunc(); err != nil {
			return err
		}
	}
	return nil
}

func (app *AgentNode) NodeInfo() string {
	data, _ := json.Marshal(app.cfg)
	return util.BytesToString(data)
}

func (app *AgentNode) Run() error {
	app.startCollect()
	return nil
}

func (app *AgentNode) Stop() {
	app.stopCollect()
	app.wg.Wait()
}
