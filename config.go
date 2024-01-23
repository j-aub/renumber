package main

import "os"

type Config struct {
	Address string
}

func LoadConfig() (config Config) {
	val, ok := os.LookupEnv("RENUMBER_ADDRESS")
	if ok {
		config.Address = val
	} else {
		config.Address = "0.0.0.0:8000"
	}

	return config
}
