package proxy

import (
	"dudu/commons/pipe"
)

func (app *ProxyNode) startPipe() error {
	app.pipePops = make([]pipe.PipePoper, 0, 10)
	if app.cfg.Proxy.HttpPipePop != nil {
		cfg := app.cfg.Proxy.HttpPipePop
		httpPipePop := pipe.NewHttpPipePop(cfg.Auth)
		// add route
		app.router.Handle(cfg.Pattern, httpPipePop)

		dataChan, err := httpPipePop.Pop()
		if err != nil {
			return err
		}

		app.wg.Add(1)
		go func() {
			defer app.wg.Done()
			app.handleDataChan(dataChan)
			app.logger.Info("close http pipe pop")
		}()
		app.pipePops = append(app.pipePops, httpPipePop)
	}
	return nil
}

func (app *ProxyNode) stopPipe() {
	for _, pipePop := range app.pipePops {
		pipePop.Stop()
	}
}

func (app *ProxyNode) handleDataChan(dataChan <-chan []byte) {
	for data := range dataChan {
		if err := app.processor.Proc(data); err != nil {
			app.logger.Warnf("proc data err:%s", err.Error())
		}
	}
}
