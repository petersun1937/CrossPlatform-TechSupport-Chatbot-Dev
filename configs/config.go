package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBString            string
	TelegramBotToken    string
	LineChannelSecret   string
	LineChannelToken    string
	ServerConfig        ServerConfig
	TelegramAPIURL      string
	TelegramWebhookURL  string
	DialogflowProjectID string
	//DBUser string
	//DBPwd  string
}

// Singleton instance of Config
var instance *Config
var once sync.Once

// Returns the singleton instance of Config
func GetConfig() *Config {
	// Ensure the config is initialized only once
	once.Do(func() {
		err := loadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
	})
	return instance
}

// Load the configuration into the singleton instance
func loadConfig() error {
	// Load the .env file only if the DATABASE_URL is not already set
	if !isEnvSet("DATABASE_URL") {
		err := godotenv.Load("configs/.env")
		if err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Initialize the config struct with environment variables
	instance = &Config{
		DBString:            os.Getenv("DATABASE_URL"),
		TelegramBotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		LineChannelSecret:   os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:    os.Getenv("LINE_CHANNEL_TOKEN"),
		TelegramAPIURL:      os.Getenv("TELEGRAM_API_URL"),
		TelegramWebhookURL:  os.Getenv("TELEGRAM_WEBHOOK_URL"),
		DialogflowProjectID: os.Getenv("DIALOGFLOW_PROJECTID"),
		ServerConfig: ServerConfig{
			Host:    os.Getenv("SERVER_HOST"),
			Port:    getEnvInt("APP_PORT", 8080),
			Timeout: getEnvDuration("SERVER_TIMEOUT", 30*time.Second),
			MaxConn: getEnvInt("SERVER_MAX_CONN", 100),
		},
	}

	// Validate required config values in a more concise way
	missingVars := []string{}
	if instance.DBString == "" {
		missingVars = append(missingVars, "DATABASE_URL")
	}
	if instance.TelegramBotToken == "" {
		missingVars = append(missingVars, "TELEGRAM_BOT_TOKEN")
	}
	if instance.TelegramAPIURL == "" {
		missingVars = append(missingVars, "TELEGRAM_API_URL")
	}

	// Return an error if any required environment variables are missing
	if len(missingVars) > 0 {
		return fmt.Errorf("required environment variables missing: %v", missingVars)
	}

	return nil
}

func (c *Config) Init() error {
	// Load environment variables
	return godotenv.Load("configs/.env")
}

// For resetting the singleton instance
func ResetConfig() {
	instance = nil     // Reset the instance for testing purposes
	once = sync.Once{} // Reset the sync.Once to allow re-initialization
}

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
