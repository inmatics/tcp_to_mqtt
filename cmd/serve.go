package cmd

import (
	"github.com/inmatics/tcp_to_mqtt/pkg/config"
	"github.com/inmatics/tcp_to_mqtt/pkg/server"
	"github.com/spf13/cobra"
	"log"
)

func extractConfig(cmd *cobra.Command) (*config.Config, error) {
	cfg := &config.Config{}

	var err error
	cfg.TcpPort, err = cmd.Flags().GetInt("tcp-port")
	if err != nil {
		return nil, err
	}

	cfg.MqttPort, err = cmd.Flags().GetInt("mqtt-port")
	if err != nil {
		return nil, err
	}

	cfg.MqttHost, err = cmd.Flags().GetString("mqtt-host")
	if err != nil {
		return nil, err
	}

	cfg.LogLevel, err = cmd.Flags().GetString("log-level")
	if err != nil {
		return nil, err
	}

	cfg.MqttUser, err = cmd.Flags().GetString("mqtt-user")
	if err != nil {
		return nil, err
	}

	cfg.MqttPassword, err = cmd.Flags().GetString("mqtt-password")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts new server",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := extractConfig(cmd)
		logFatal(err)

		server.Start(config)
	},
}

func init() {
	serveCmd.Flags().IntP("tcp-port", "p", 3064, "TCP port to listen on")
	serveCmd.Flags().IntP("mqtt-port", "m", 1883, "MQTT broker port to publish to")
	serveCmd.Flags().StringP("mqtt-host", "i", "localhost", "MQTT broker host to publish to")

	serveCmd.Flags().StringP("mqtt-user", "u", "", "MQTT broker user")
	serveCmd.Flags().StringP("mqtt-password", "P", "", "MQTT broker password")

	rootCmd.PersistentFlags().StringP("log-level", "l", "", "Zone ID")
	rootCmd.AddCommand(serveCmd)
}
