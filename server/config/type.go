package config

// Config for all app
type Config struct {
	Nats   NatsConfig   `json:"nats"`
	Server ServerConfig `json:"server"`
}

// information about nats for connection to it
type NatsConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// requirements for start server  
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
