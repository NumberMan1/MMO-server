package bootstrap

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	cfgFile string
	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "game_server",
		Short: "Game Server",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func init() {
	// Bind the config flag
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	if err := RootCmd.MarkPersistentFlagRequired("config"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}

	// Initialize viper to read the config file
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".config" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".config")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		_, _ = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
