package server

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/inmatics/tcp_to_mqtt/teltonika"
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
	messages := make(chan string)

	go listen(messages, client, logger)()

	l, err := net.Listen("tcp", "0.0.0.0:"+connPort)
	logFatal(err)
	defer l.Close()

	fmt.Println("TCP server listening on port " + connPort)
	for {
		conn, err := l.Accept()
		logFatal(err)
		go teltonika.HandleRequest(conn, messages)
	}
}

func listen(messages chan string, client mqtt.Client, logger *slog.Logger) func() {
	return func() {
		for msg := range messages {
			token := client.Publish("devices/new", 0, false, msg)
			logger.Debug("new message",
				slog.String("msg", msg),
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
