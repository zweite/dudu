package main

import (
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

	rootCmd.AddCommand(commands.NewAgentNodeCmd(agent.NewAgentNode))

	cmd := cli.PrepareBaseCmd(rootCmd,
		"TM",
		config.DefaultDir,
	)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
