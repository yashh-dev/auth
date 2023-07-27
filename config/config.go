package config

import (
	"github.com/caarlos0/env/v9"
	"log"
)

type Config struct {
	JWTSecret string `env:"JWT_SECRET,required"`
	PEPPER    string `env:"PEPPER,required"`
	DBUrl     string `env:"DB_URL,required"`
	RedisHOST string `env:"REDIS_HOST,required"`
	RedisPass string `env:"REDIS_PASS"`
	RabbitMQ  string `env:"RABBITMQ,required"`
}

func GetConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
	return &cfg
}
