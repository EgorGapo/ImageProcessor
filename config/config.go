package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	// Database
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	// Redis
	RedisHost string
	RedisPort string

	// RabbitMQ
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string

	// App
	AppHost string
	AppPort string
	AppAddr string
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = ".env"
	}

	envMap := make(map[string]string)

	// Load from file if exists
	if file, err := os.Open(path); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				envMap[key] = value
			}
		}
	}

	// Override with environment variables
	for key := range envMap {
		if val, ok := os.LookupEnv(key); ok && val != "" {
			envMap[key] = val
		}
	}

	// Helper function to get env value
	getEnv := func(key, defaultVal string) string {
		if val, ok := envMap[key]; ok {
			return val
		}
		return defaultVal
	}

	cfg := &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "db"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "12345"),
		PostgresDB:       getEnv("POSTGRES_DB", "postgres"),
		RedisHost:        getEnv("REDIS_HOST", "redis"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RabbitMQHost:     getEnv("RABBITMQ_HOST", "broker"),
		RabbitMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQUser:     getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPassword: getEnv("RABBITMQ_PASSWORD", "guest"),
		AppHost:          getEnv("APP_HOST", "app"),
		AppPort:          getEnv("APP_PORT", "8080"),
		AppAddr:          getEnv("APP_ADDR", ":8080"),
	}

	return cfg, nil
}

func (c *Config) PostgresConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c *Config) RabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.RabbitMQUser, c.RabbitMQPassword, c.RabbitMQHost, c.RabbitMQPort)
}

func (c *Config) AppURL() string {
	return fmt.Sprintf("http://%s:%s", c.AppHost, c.AppPort)
}
