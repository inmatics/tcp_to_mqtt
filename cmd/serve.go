package cmd

import (
	"github.com/inmatics/tcp_to_mqtt/pkg/server"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts new server",
	Run: func(cmd *cobra.Command, args []string) {

		tcpPort, err := cmd.Flags().GetInt("tcp-port")
		logFatal(err)

		mqttPort, err := cmd.Flags().GetInt("mqtt-port")
		logFatal(err)

		mqttHost, err := cmd.Flags().GetString("mqtt-host")
		logFatal(err)

		server.Start(strconv.Itoa(tcpPort), mqttHost, strconv.Itoa(mqttPort))
	},
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	serveCmd.Flags().IntP("tcp-port", "p", 3064, "TCP port to listen on")
	serveCmd.Flags().IntP("mqtt-port", "m", 1883, "MQTT broker host to publish to")
	serveCmd.Flags().StringP("mqtt-host", "i", "localhost", "MQTT broker port to publish to")

	rootCmd.AddCommand(serveCmd)
}
