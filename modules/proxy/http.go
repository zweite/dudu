package proxy

import (
	"net"
	"net/http"
)

func (app *ProxyNode) initHttpHandle() (err error) {
	return
}

func (app *ProxyNode) startHttpServ() (err error) {
	listener, err := net.Listen("tcp", app.cfg.Proxy.HttpAddr)
	if err != nil {
		return
	}

	app.logger.Infof("start http serv on %s", app.cfg.Proxy.HttpAddr)

	app.listener = listener
	return http.Serve(listener, app.router)
}

func (app *ProxyNode) stopHttpServ() {
	app.listener.Close()
	app.logger.Info("close http serv")
}
