package command

import (
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar service",
	Long:  "A simple calendar service with HTTP API",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to config file (required)")
	rootCmd.MarkPersistentFlagRequired("config")
}

func Execute() error {
	return rootCmd.Execute()
}
