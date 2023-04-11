package main

import (
	"encoding/hex"
	"github.com/inmatics/tcp_to_mqtt/pkg/streams"
	"log"
	"net"
)

const (
	ExampleImei1    = "000f383639373338303636303632343737"
	ExampleAvlData1 = "0000000000000044080100000186e9f5bf1000dd1656fdeb70b8ea001700290e0000000c05ef00f0001503c800450105b5000ab60007422fc143101744000002f1000b05861000638e0b00010000f41f"
)

func main() {
	ConnPort := "3064"
	connection, _ := net.Dial("tcp", "localhost"+":"+ConnPort)
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Println("Error closing connection.")
		}
	}(connection)

	decodeString, err := hex.DecodeString(ExampleImei1)
	if err != nil {
		return
	}
	_, err = connection.Write(decodeString)
	if err != nil {
		logFatal(err)
	}

	buf := make([]byte, 1024)
	n, err := connection.Read(buf)
	logFatal(err)

	streams.ToInt8(buf[:n])
	decodeString, err = hex.DecodeString(ExampleAvlData1)
	logFatal(err)

	_, err = connection.Write(decodeString)
	logFatal(err)

	_, err = connection.Read(buf)
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
