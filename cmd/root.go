package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
)

// NewRootCmd returns the root command for relayer.
func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "subnode",
		Short: "subnode is a smart reverse proxy server for cosmos based chains",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// reads `homeDir/config.yaml` into `var config *Config` before each command
		// if err := initConfig(rootCmd); err != nil {
		// 	return err
		// }

		return nil
	}

	// --app flag
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "subnode.yaml", "path to the config file")
	if err := viper.BindPFlag("config-file", rootCmd.PersistentFlags().Lookup("config-file")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(
		startCmd(),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
