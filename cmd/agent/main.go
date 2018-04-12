package main

import (
	"os"

	"dudu/cmd/agent/commands"
	"dudu/cmd/version"
	"dudu/commons/cli"
	"dudu/config"
	"dudu/modules/agent"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		version.VersionCmd)

	// new agent node func
	nodeFunc := agent.NewAgentNode

	rootCmd.AddCommand(commands.NewAgentNodeCmd(nodeFunc))

	cmd := cli.PrepareBaseCmd(rootCmd,
		"TM",
		os.ExpandEnv(config.DefaultDir),
	)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
