package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/inmatics/tcp_to_mqtt/pkg/config"
	"github.com/inmatics/tcp_to_mqtt/pkg/teltonika"
	"golang.org/x/exp/slog"
)

func Start(cfg *config.Config) {
	logger := getLogger(cfg.LogLevel)

	mqttPort := strconv.Itoa(cfg.MqttPort)
	server := cfg.MqttHost + ":" + mqttPort

	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.SetUsername(cfg.MqttUser)
	opts.SetPassword(cfg.MqttPassword)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Println("Error connecting to broker " + server)
		logFatal(token.Error())
	}
	messages := make(chan teltonika.Record)

	go listen(messages, client, logger)()

	tcpPort := strconv.Itoa(cfg.TcpPort)
	l, err := net.Listen("tcp", "0.0.0.0:"+tcpPort)
	logFatal(err)
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Println("error closing connection")
		}
	}(l)

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

func listen(recordsChannel <-chan teltonika.Record, client mqtt.Client, logger *slog.Logger) func() {
	return func() {
		for record := range recordsChannel {
			bytes, err := json.Marshal(record)
			if err != nil {
				logger.Error("error marshalling teltonika record: ", record, err.Error())
				return
			}

			client.Publish("devices/new", 1, false, string(bytes))
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
