package cmd

import (
	"os"

	"github.com/parsaaes/jasoos-telegram-bot/cmd/server"
	"github.com/parsaaes/jasoos-telegram-bot/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Execute() {
	cmd := &cobra.Command{
		Use:   "jasoos-telegram-bot",
		Short: "jasoos telegram bot",
	}

	cfg := config.New()

	cmd.AddCommand(server.Cmd(cfg))

	if err := cmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
