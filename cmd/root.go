package cmd

import (
	"github.com/parsaaes/jasoos-telegram-bot/cmd/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	cmd := &cobra.Command{
		Use:   "jasoos-telegram-bot",
		Short: "jasoos telegram bot server",
	}

	cmd.AddCommand(server.Cmd())

	if err := cmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
