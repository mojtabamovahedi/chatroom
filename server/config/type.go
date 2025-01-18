package config

type Config struct {
	Nats NatsConfig `json:"nats"`
}

type NatsConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Subject string `json:"subject"`
}
