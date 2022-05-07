package config

import (
	"blockchain/pkg/config/logger"
	"blockchain/pkg/config/system"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

var config *Configuration
var once sync.Once

type Configuration struct {
	System system.Config
	Logger logger.Config
}

func LoadConfig() {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Print(fmt.Sprintf("not found .env file: %v", err))
		}
		if os.Getenv("ENV") == "" {
			log.Fatal("env loading error. ")
		}

		config = &Configuration{
			System: system.NewConfig(),
			Logger: logger.NewConfig(),
		}
	})
}

func GetConfig() *Configuration {
	return config
}
