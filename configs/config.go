package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBString string
	//DBUser string
	//DBPwd  string
}

func NewConfig(dbURL string) *Config {
	return &Config{
		DBString: dbURL,
	}
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

// GetDBString returns the full connection string
func (c *Config) GetDBString() string {
	return c.DBString
}

type ServerConfig struct {
	Host    string
	Port    int // generally int
	Timeout time.Duration
	MaxConn int
}
