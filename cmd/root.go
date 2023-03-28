package cmd

import (
	"fmt"
	"github.com/notional-labs/subnode/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfg  *config.Config
	conf string // path to config file
)

// NewRootCmd returns the root command
func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "subnode",
		Short: "subnode is a smart reverse proxy server for cosmos based chains",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		cfg, err := config.LoadConfigFromFile(conf)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", cfg)

		return nil
	}

	// --config flag
	rootCmd.PersistentFlags().StringVar(&conf, "conf", "subnode.yaml", "path to the config file")
	if err := viper.BindPFlag("conf", rootCmd.PersistentFlags().Lookup("conf")); err != nil {
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

func GetConfig() *config.Config {
	return cfg
}
