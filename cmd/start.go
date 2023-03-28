package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start subnode server",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("startCmd")

			return nil
		},
	}
	return cmd
}
