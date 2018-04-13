package commands

import (
	"fmt"

	"dudu/commons/event"
	"dudu/config"
	"dudu/modules/agent"

	"github.com/spf13/cobra"
)

func AddNodeFlags(cmd *cobra.Command) {
	cmd.Flags().String(config.AddrFlag, cfg.Agent.Addr, "Node addr")
	cmd.Flags().Bool(config.DebugFlag, cfg.Agent.Debug, "Node debug")
}

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

	AddNodeFlags(cmd)
	return cmd
}
