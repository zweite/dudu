package agent

import "dudu/config"

// agent node context

type AgentContext struct{}

func NewAgentContext(cfg *config.AgentConfig) *AgentContext {
	return &AgentContext{}
}
