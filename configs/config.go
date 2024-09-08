package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBString          string
	TelegramBotToken  string
	LineChannelSecret string
	LineChannelToken  string
	ServerConfig      ServerConfig
	//DBUser string
	//DBPwd  string
}

/*func NewConfig(dbURL string) *Config {
	return &Config{
		DBString: dbURL,
	}
}*/

// NewConfig loads environment variables and initializes the config
func NewConfig() (*Config, error) {
	// Only load .env if required environment variables are not already set
	if !isEnvSet("DATABASE_URL") || !isEnvSet("TELEGRAM_BOT_TOKEN") || !isEnvSet("LINE_CHANNEL_SECRET") || !isEnvSet("LINE_CHANNEL_TOKEN") {
		err := godotenv.Load("configs/.env")
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// // Load environment variables from .env file
	// err := godotenv.Load("configs/.env")
	// if err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

	// Initialize the config struct with environment variables
	config := &Config{
		DBString:          os.Getenv("DATABASE_URL"),
		TelegramBotToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:  os.Getenv("LINE_CHANNEL_TOKEN"),
		ServerConfig: ServerConfig{
			Host:    os.Getenv("SERVER_HOST"),
			Port:    getEnvInt("APP_PORT", 8080),
			Timeout: getEnvDuration("SERVER_TIMEOUT", 30*time.Second),
			MaxConn: getEnvInt("SERVER_MAX_CONN", 100),
		},
	}

	// Validate required config values
	if config.DBString == "" || config.TelegramBotToken == "" {
		return nil, fmt.Errorf("required environment variables missing")
	}

	return config, nil
}

func (c *Config) Init() error {
	// Load environment variables
	return godotenv.Load("configs/.env")
}

/*
func (c *Config) GetTelegramBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func (c *Config) GetLineSecret() string {
	return os.Getenv("LINE_CHANNEL_SECRET")
}

func (c *Config) GetLineToken() string {
	return os.Getenv("LINE_CHANNEL_TOKEN")
}

// GetDBString returns the full connection string
func (c *Config) GetDBString() string {
	return c.DBString
}*/

// Helper function to check if an environment variable is set
func isEnvSet(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// Utility function to get environment variable as an integer
func getEnvInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}

// Utility function to get environment variable as a duration
func getEnvDuration(name string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(name); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultVal
}

type ServerConfig struct {
	Host    string
	Port    int // generally int
	Timeout time.Duration
	MaxConn int
}
