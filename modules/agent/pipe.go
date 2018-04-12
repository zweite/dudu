package agent

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"dudu/config"
)

// 日志传输管道
type Pipe interface {
	Push([]byte) error
	// Pop() ([]byte, error)
}

func (app *AgentNode) initPipe() error {
	switch app.cfg.Agent.Pipe {
	case "http":
		app.pipe = NewHttpPipe(app.cfg.Agent.HttpPipe)
	case "":
		return errors.New("pipe can't be empty")
	default:
		return fmt.Errorf("not found pipe 【%s】", app.cfg.Agent.Pipe)
	}
	return nil
}

type HttpPipe struct {
	auth   string // 认证信息
	addr   string // 传输地址
	client *http.Client
}

func NewHttpPipe(httpPipeCfg *config.HttpPipeConfig) *HttpPipe {
	return &HttpPipe{
		auth:   httpPipeCfg.Auth,
		addr:   httpPipeCfg.Addr,
		client: new(http.Client),
	}
}

func (h *HttpPipe) Push(data []byte) (err error) {
	req, err := http.NewRequest("POST", h.addr, bytes.NewReader(data))
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Basic "+h.auth)
	resp, err := h.client.Do(req)
	if err != nil {
		return
	}

	resp.Body.Close()
	return
}
