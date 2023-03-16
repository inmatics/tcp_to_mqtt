package cmd

import (
	"github.com/inmatics/tcp_to_mqtt/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts new server",
	Run: func(cmd *cobra.Command, args []string) {
		ConnPort := "3064"
		mqttBrokerHost := "localhost"
		mqttBrokerPort := "1883"
		server.Start(ConnPort, mqttBrokerHost, mqttBrokerPort)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
