package proxy

import "dudu/config"

type Forward struct {
	cfg    *config.ProxyForwardConfig
	resCfg *config.ResourceConfig
}

func NewForward(cfg *config.ProxyForwardConfig, resCfg *config.ResourceConfig) (*Forward, error) {
	return &Forward{
		cfg:    cfg,
		resCfg: resCfg,
	}, nil
}

// 转发无需解压
func (f *Forward) Proc(data []byte) (err error) {
	return
}

func (f *Forward) Stop() {

}
