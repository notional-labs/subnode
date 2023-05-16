package cmd

import (
	"github.com/notional-labs/subnode/pkg/server"
	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start subnode server",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			server.Start()

			return nil
		},
	}
	return cmd
}
