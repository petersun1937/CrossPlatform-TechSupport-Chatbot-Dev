package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	// Load environment variables
	return godotenv.Load("configs/.env")
}

func (c *Config) GetTelegramBotToken() string {
	return os.Getenv("TELEGRAM_BOT_TOKEN")
}

func (c *Config) GetLineSecret() string {
	return os.Getenv("LINE_CHANNEL_SECRET")
}

func (c *Config) GetLineToken() string {
	return os.Getenv("LINE_CHANNEL_TOKEN")
}
