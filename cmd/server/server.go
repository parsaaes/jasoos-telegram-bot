package server

import (
	"github.com/parsaaes/jasoos-telegram-bot/config"
	"github.com/spf13/cobra"
)

// Cmd create server command
func Cmd(cfg config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start the game server",
		Run: func(cmd *cobra.Command, args []string) {
			panic("implement me")
		},
	}

	return cmd
}
