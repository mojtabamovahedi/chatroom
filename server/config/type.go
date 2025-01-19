package config

type Config struct {
	Nats   NatsConfig   `json:"nats"`
	Server ServerConfig `json:"server"`
}

type NatsConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Subject string `json:"subject"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
