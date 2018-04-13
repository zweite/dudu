package agent

import (
	"errors"
	"fmt"

	"dudu/commons/pipe"
	"dudu/config"
)

func (app *AgentNode) initPipe() error {
	switch app.cfg.Agent.Pipe {
	case "http":
		app.pipe = NewHttpPipe(app.cfg.Agent.HttpPipePush)
	case "":
		return errors.New("pipe can't be empty")
	default:
		return fmt.Errorf("not found pipe 【%s】", app.cfg.Agent.Pipe)
	}
	return nil
}

func NewHttpPipe(httpPipeCfg *config.HttpPipePushConfig) pipe.Pipe {
	return pipe.NewHttpPipe(
		pipe.NewHttpPipePush(httpPipeCfg.Addr, httpPipeCfg.Auth),
		nil,
	)
}
