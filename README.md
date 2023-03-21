# TCP to MQTT
## Introduction
tcp_to_mqtt is a project to collect TCP packets from iOT devices that transmit on TCP and relay them to a MQTT broker.

## Installation
You can build the project running
```shell
make build
```

or you can install it to your GOPATH using 
```shell
go install
```

## Usage
```shell
tcp_to_mqtt serve --mqtt-port 9999 --log-level debug
```
 

