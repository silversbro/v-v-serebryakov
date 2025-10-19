package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConfig   `toml:"logger"`
	Database DatabaseConfig `toml:"database"`
	Cache    CacheConfig    `toml:"cache"`
	Queue    QueueConfig    `toml:"queue"`
	Redis    RedisConfig    `toml:"redis"`
	RabbitMQ RabbitMQConfig `toml:"rabbitmq"`
	Service  ServiceConfig  `toml:"service"`
	HTTP     HTTPConfig     `toml:"http"`
}

type LoggerConfig struct {
	Driver string `toml:"driver"` // stdout, file, sentry
	Type   string `toml:"type"`   // daily, single
	Level  string `toml:"level"`  // debug, info, notice, error, alert
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	SSLMode  string `toml:"sslmode"`
}

type CacheConfig struct {
	Driver  string `toml:"driver"` // redis, memory
	Address string `toml:"address"`
}

type QueueConfig struct {
	Driver  string `toml:"driver"` // rabbitmq, memory
	Address string `toml:"address"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type RabbitMQConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	VHost    string `toml:"vhost"`
}

type ServiceConfig struct {
	Name        string `toml:"name"`
	Environment string `toml:"environment"` // development, staging, production
	Debug       bool   `toml:"debug"`
}

type HTTPConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}
