package main

import (
	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start subnode server",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			//Init()
			//StartRpcServer()
			//StartApiServer()
			StartGrpcServer()

			select {}

			return nil
		},
	}
	return cmd
}
