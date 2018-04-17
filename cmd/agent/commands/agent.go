package commands

import (
	"fmt"

	"dudu/commons/event"
	"dudu/modules/agent"

	"github.com/spf13/cobra"
)

func AddAgentNodeFlags(cmd *cobra.Command) {}

func NewAgentNodeCmd(nodeProvider agent.AgentNodeProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run the agent node",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create & Run node
			n, err := nodeProvider(cfg, logger)
			if err != nil {
				return fmt.Errorf("Failed to create node: %v", err)
			}

			if err := n.Init(); err != nil {
				return fmt.Errorf("Failed to init node: %v", err)
			} else {
				logger.Info("Inited node", "nodeInfo", n.NodeInfo())
			}
			// Trap signal, run forever.
			go n.Run()
			event.AddHook(event.HightPriority, n.Stop)
			event.WaitExit()
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("will be exit")
			return nil
		},
	}

	AddAgentNodeFlags(cmd)
	return cmd
}
