package server

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/inmatics/tcp_to_mqtt/pkg/teltonika"
	"golang.org/x/exp/slog"
	"log"
	"net"
	"os"
)

func Start(connPort string, mqttBrokerHost string, mqttBrokerPort string, logLevel string) {
	logger := getLogger(logLevel)

	opts := mqtt.NewClientOptions().AddBroker(mqttBrokerHost + ":" + mqttBrokerPort)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		logFatal(token.Error())
	}
	messages := make(chan teltonika.Record)

	go listen(messages, client, logger)()

	l, err := net.Listen("tcp", "0.0.0.0:"+connPort)
	logFatal(err)
	defer l.Close()

	fmt.Println("TCP server listening on port " + connPort)
	fmt.Println("Relaying MQTT messages to " + mqttBrokerHost + " on port " + mqttBrokerPort)
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
