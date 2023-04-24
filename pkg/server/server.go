package server

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/inmatics/tcp_to_mqtt/pkg/config"
	"github.com/inmatics/tcp_to_mqtt/pkg/teltonika"
	"golang.org/x/exp/slog"
	"log"
	"net"
	"os"
	"strconv"
)

func Start(cfg *config.Config) {
	logger := getLogger(cfg.LogLevel)

	mqttPort := strconv.Itoa(cfg.MqttPort)
	opts := mqtt.NewClientOptions().AddBroker(cfg.MqttHost + ":" + mqttPort)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		logFatal(token.Error())
	}
	messages := make(chan teltonika.Record)

	go listen(messages, client, logger)()

	tcpPort := strconv.Itoa(cfg.TcpPort)
	l, err := net.Listen("tcp", "0.0.0.0:"+tcpPort)
	logFatal(err)
	defer l.Close()

	fmt.Println("TCP server listening on port " + tcpPort)
	fmt.Println("Relaying MQTT messages to " + cfg.MqttHost + " on port " + mqttPort)
	for {
		conn, err := l.Accept()
		logFatal(err)
		go teltonika.HandleRequest(conn, messages, logger)
	}
}

func getLogger(level string) *slog.Logger {
	logLevel := new(slog.LevelVar)
	logger := slog.New(slog.HandlerOptions{Level: logLevel}.NewJSONHandler(os.Stderr))
	if level == "debug" {
		logLevel.Set(slog.LevelDebug)
	}
	return logger
}

func listen(records chan teltonika.Record, client mqtt.Client, logger *slog.Logger) func() {
	return func() {
		for record := range records {
			bytes, err := json.Marshal(record)
			logFatal(err)

			client.Publish("devices/new", 0, false, string(bytes))
			client.Publish("devices/"+record.Imei, 0, false, string(bytes))
			logger.Debug("new message for imei: "+record.Imei,
				slog.String("msg", string(bytes)),
			)

		}
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
