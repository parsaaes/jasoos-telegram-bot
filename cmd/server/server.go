package server

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start the game server",
		Run: func(cmd *cobra.Command, args []string) {
			panic("implement me")
		},
	}

	return cmd
}
