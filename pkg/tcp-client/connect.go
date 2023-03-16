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
	connection, err := net.Dial("tcp", "localhost"+":"+ConnPort)
	defer connection.Close()

	decodeString, err := hex.DecodeString(ExampleImei1)
	if err != nil {
		return
	}
	_, err = connection.Write(decodeString)

	buf := make([]byte, 1024)
	n, err := connection.Read(buf)

	streams.ToInt8(buf[:n])
	decodeString, err = hex.DecodeString(ExampleAvlData1)

	_, err = connection.Write(decodeString)
	n, err = connection.Read(buf)
	if err != nil {
		log.Fatal("Error reading")
	}
}
