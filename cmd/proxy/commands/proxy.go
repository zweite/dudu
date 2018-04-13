package commands

import (
	"fmt"

	"dudu/commons/event"
	"dudu/config"
	"dudu/modules/proxy"

	"github.com/spf13/cobra"
)

func AddProxyNodeFlags(cmd *cobra.Command) {
	cmd.Flags().String(config.HttpAddrFlag, cfg.Proxy.HttpAddr, "node http addr")
	cmd.Flags().String(config.ProxyModeFlag, cfg.Proxy.Mode, "proxy mode")
}

func NewProxyNodeCmd(nodeProvider proxy.ProxyNodeProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Run the proxy node",
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
	}

	AddProxyNodeFlags(cmd)
	return cmd
}
