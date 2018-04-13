package agent

import (
	"fmt"

	"dudu/commons/compactor"
)

func (app *AgentNode) initCompactor() error {
	switch app.cfg.Agent.Compactor {
	case "snappy":
		app.compactor = compactor.NewSnappy()
	case "gzip":
		app.compactor = compactor.NewGZip()
	case "":
		app.compactor = nil
	default:
		return fmt.Errorf("not found compactor 【%s】", app.cfg.Agent.Compactor)
	}
	return nil
}
