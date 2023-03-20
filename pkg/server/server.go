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

func Start(connPort string, mqttBrokerHost string, mqttBrokerPort string) {
	textHandler := slog.NewTextHandler(os.Stdout)
	logger := slog.New(textHandler)
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
		go teltonika.HandleRequest(conn, messages)
	}
}

func listen(messages chan teltonika.Record, client mqtt.Client, logger *slog.Logger) func() {
	return func() {
		for msg := range messages {
			topic := "devices/new"
			bytes, err := json.Marshal(msg)
			logFatal(err)

			token := client.Publish(topic, 0, false, msg)
			logger.Debug("new message",
				slog.String("msg", string(bytes)),
			)
			token.Wait()
		}
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
