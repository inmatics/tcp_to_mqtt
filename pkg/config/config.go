package config

type Config struct {
	TcpPort      int
	MqttPort     int
	MqttHost     string
	LogLevel     string
	MqttPassword string
	MqttUser     string
}
