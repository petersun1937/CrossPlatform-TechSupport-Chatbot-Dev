package bot

import (
	"Tg_chatbot/config"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBot interface {
	Run() error
}

type lineBot struct {
	secret     string
	token      string
	lineClient *linebot.Client
}

func NewLineBot(conf *config.Config) Bot {
	return &lineBot{
		secret: conf.GetLineSecret(),
		token:  conf.GetLineToken(),
	}
}

func (b *lineBot) Run() error {
	// Initialize Linebot
	lineClient, err := linebot.New(b.secret, b.token) // create new BotAPI instance using the channel token and secret
	if err != nil {
		return err
	}

	b.lineClient = lineClient
	return nil
}
