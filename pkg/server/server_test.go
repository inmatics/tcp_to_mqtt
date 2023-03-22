package server

import (
	"encoding/hex"
	"github.com/inmatics/tcp_to_mqtt/pkg/streams"
	"log"
	"net"
	"testing"
	"time"
)

const (
	// ExampleImei1 Info for Teltonika codec 8 -> https://wiki.teltonika-gps.com/view/Codec
	ExampleImei1    = "000f383636333831303531333132313630"
	ExampleAvlData1 = "000000000000004308020000016B40D57B480100000000000000000000000000000001010101000000000000016B40D5C198010000000000000000000000000000000101010101000000020000252C"
)

func TestStart(t *testing.T) {
	ConnPort := "3064"
	mqttBrokerHost := "tcp://mqtt.inmatics.io"
	mqttBrokerPort := "9001"
	go Start(ConnPort, mqttBrokerHost, mqttBrokerPort, "debug")
	time.Sleep(500 * time.Millisecond)

	connection, err := net.Dial("tcp", "localhost"+":"+ConnPort)
	logFatal(err)
	defer connection.Close()

	// First, when module connects to server, module sends its IMEI. First comes short identifying number of bytes written and then goes IMEI as text (bytes).
	// For example, IMEI 356307042441013 would be sent as 000F333536333037303432343431303133.
	// First two bytes denote IMEI length. In this case 0x000F means, that IMEI is 15 bytes long.

	// Then module starts to send first AVL data packet. After server receives packet and parses it, server must report to module number of data received as integer (four bytes).
	// If sent data number and reported by server doesn’t match module resends sent data.
	decodeString, err := hex.DecodeString(ExampleImei1)
	if err != nil {
		return
	}
	_, err = connection.Write(decodeString)
	logFatal(err)

	// Receive data from the server
	buf := make([]byte, 1024)

	n, err := connection.Read(buf)
	logFatal(err)

	serverResponse := streams.ToInt8(buf[:n])

	// After receiving IMEI, server should determine if it would accept data from this module.
	asserServerResponse(serverResponse)
	decodeString, err = hex.DecodeString(ExampleAvlData1)
	logFatal(err)

	_, err = connection.Write(decodeString)
	logFatal(err)

	n, err = connection.Read(buf)
	logFatal(err)

	avlServerResponse := streams.ToInt32(buf[:n])
	if avlServerResponse != 2 {
		t.Error("Server response should be 2")
	}

	//os.Exit(0)
}

func asserServerResponse(serverResponse int8) {
	// If the server accepted the data, it should reply 0x01
	if serverResponse != 1 {
		log.Fatal("Server response is not 1")
	}
}
