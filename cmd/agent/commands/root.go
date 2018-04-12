package commands

import (
	"os"

	"dudu/cmd/version"
	"dudu/commons/log"
	"dudu/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg    = config.DefaultConfig()
	logger = log.NewLogger(os.Stdout, config.DefaultLogLevelInt())
)

func ParseConfig() (*config.Config, error) {
	conf := config.ParseConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.SetRoot(conf.RootDir)
	return conf, err
}

var RootCmd = &cobra.Command{
	Use:   "agent",
	Short: "agent in Go",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == version.VersionCmd.Name() {
			return nil
		}

		cfg, err = ParseConfig()
		if err != nil {
			return err
		}

		logger = log.ParseLogLevel(cfg.LogLevel, config.DefaultLogLevel(), logger)
		return nil
	},
}
