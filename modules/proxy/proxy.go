package proxy

import (
	"dudu/commons/pipe"
	"encoding/json"
	"net"
	"sync"

	"github.com/gorilla/mux"

	"dudu/commons/log"
	"dudu/commons/util"
	"dudu/config"
)

// 代理节点

type ProxyNodeProvider func(*config.Config, log.Logger) (*ProxyNode, error)

type ProxyNode struct {
	router    *mux.Router
	wg        *sync.WaitGroup
	cfg       *config.Config
	pipePops  []pipe.PipePoper
	processor Processor // 接收数据处理器
	listener  net.Listener
	logger    log.Logger
}

func NewProxyNode(cfg *config.Config, logger log.Logger) (*ProxyNode, error) {
	return &ProxyNode{
		router: mux.NewRouter(),
		wg:     new(sync.WaitGroup),
		cfg:    cfg,
		logger: logger,
	}, nil
}

// 初始化
func (app *ProxyNode) Init() (err error) {
	if err = app.initHttpHandle(); err != nil {
		return
	}

	if err = app.initProc(app.cfg.Proxy.Mode); err != nil {
		return
	}
	return nil
}

// 节点信息
func (app *ProxyNode) NodeInfo() string {
	data, _ := json.Marshal(app.cfg)
	return util.BytesToString(data)
}

// 运行节点
func (app *ProxyNode) Run() {
	if err := app.run(); err != nil {
		panic(err)
	}
}

func (app *ProxyNode) run() (err error) {
	if err = app.startPipe(); err != nil {
		return
	}

	if err = app.startHttpServ(); err != nil {
		return
	}

	return nil
}

// 停止节点
func (app *ProxyNode) Stop() {
	app.stopHttpServ()
	app.stopPipe()
	app.stopProc()
	app.wg.Wait()
	app.logger.Info("will close proxy node")
}
