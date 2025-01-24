package config

type Config struct {
	Nats   NatsConfig   `json:"nats"`
	Server ServerConfig `json:"server"`
}

type NatsConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
