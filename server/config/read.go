package config

import (
	"encoding/json"
	"os"
)

// Read config from path 
func ReadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return config, json.Unmarshal(data, &config)
}


// Must read config from path if there was an error it will panic
func MustReadConfig(path string) Config {
	config, err := ReadConfig(path)
	if err != nil {
		panic(err)
	}
	return config
}
