package logger

import "os"

type Config struct {
	LogEncoding string
	LogLevel    string
}

func NewConfig() Config {
	return Config{
		LogEncoding: os.Getenv("LOG_ENCODING"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}
}
