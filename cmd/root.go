// Package cmd entry point for application commands
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tcp_to_mqtt",
	Short: "Tool to create a daemon that list for TCP connections and publishes MQTT topics",
}

// Execute root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
