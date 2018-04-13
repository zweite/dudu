package proxy

import "fmt"

type Processor interface {
	Proc([]byte) error
}

// 目前支持两种模式，一种是forward，另外一种是persistence
func (app *ProxyNode) initProc(mode string) (err error) {
	switch mode {
	case "":
		fallthrough
	case "persistence":
		app.processor, err = NewPersistence(app.cfg.Proxy.Persistence, app.cfg.Resource, app.logger)
	case "forward":
		app.processor, err = NewForward(app.cfg.Proxy.Forward, app.cfg.Resource)
	default:
		return fmt.Errorf("can't support 【%s】 model", mode)
	}
	return nil
}
