package main

import (
	"dudu/cmd/proxy/commands"
	"dudu/cmd/version"
	"dudu/commons/cli"
	"dudu/config"
	"dudu/modules/proxy"
)

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		version.VersionCmd)

	rootCmd.AddCommand(commands.NewProxyNodeCmd(proxy.NewProxyNode))

	cmd := cli.PrepareBaseCmd(rootCmd,
		"TM",
		config.DefaultDir,
	)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
