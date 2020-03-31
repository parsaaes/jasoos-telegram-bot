package server

import (
	"github.com/parsaaes/jasoos-telegram-bot/config"
	"github.com/parsaaes/jasoos-telegram-bot/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Cmd create server command
func Cmd(cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start the game server",
		Run: func(cmd *cobra.Command, args []string) {
			eng, err := engine.New(cfg.Token)
			if err != nil {
				logrus.Fatalf("server: cannot create engine: %s", err.Error())
				return
			}

			eng.Run()
		},
	}

	return cmd
}
