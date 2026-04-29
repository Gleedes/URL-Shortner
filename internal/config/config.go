package config

import (
	"fmt"
	"os"

	"github.com/subosito/gotenv"
)

type config struct {
	ServerPort   string
	DatabasePath string
}

func Load() (*config, error) {
	if err := gotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}
	Config := &config{
		ServerPort:   os.Getenv("SERVER_PORT"),
		DatabasePath: os.Getenv("DATABASE_PATH"),
	}

	if Config.ServerPort == "" {
		return nil, fmt.Errorf("SERVER_PORT is not set")
	}

	return Config, nil
}
