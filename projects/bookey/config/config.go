package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

// Config ...
type Config struct {
	Debug  bool `env:"DEBUG,required"`
	Server serverConfig
	Db     dbConfig
}

// serverConfig ..
type serverConfig struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

// dbConfig
type dbConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
}

// AppConfig ...
func AppConfig() *Config {
	var config Config

	if err := envdecode.StrictDecode(&config); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &config
}
