package server

import (
	"Tg_chatbot/bot"
	config "Tg_chatbot/configs"
	"Tg_chatbot/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config  *config.Config
	Service *service.Service // for database operations
	Router  *gin.Engine
	Bots    map[string]bot.Bot
	// LineBot    bot.LineBot
	// TgBot      bot.TgBot
	// FbBot      bot.FbBot
	// GeneralBot bot.GeneralBot // custom frontend
}

func (a App) Run() error {
	// initialize http server routes from app struct
	go a.RunRoutes(a.Config, a.Service)

	// running bots
	for _, bot := range a.Bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}

	return nil
}

func NewApp(conf *config.Config, srv *service.Service, bots map[string]bot.Bot) *App {
	// // Initialize the line bot
	// lineBot, err := bot.NewLineBot(conf, srv)
	// if err != nil {
	// 	log.Fatal("Failed to initialize LINE bot:", err)
	// }

	// // Initialize the tg bot
	// tgBot, err := bot.NewTGBot(conf, srv)
	// if err != nil {
	// 	//log.Fatal("Failed to initialize Telegram bot:", err)
	// 	fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	// }

	// // Initialize the tg bot
	// fbBot, err := bot.NewFBBot(conf, srv)
	// if err != nil {
	// 	fmt.Printf("Failed to initialize Messenger bot: %s", err.Error())
	// }

	// // Initialize the general bot
	// generalBot := bot.NewGeneralBot(nil)

	return &App{
		Config:  conf,
		Service: srv,
		Router:  gin.Default(),
		Bots:    bots,
		// LineBot:    lineBot,    // Store the initialized LineBot
		// TgBot:      tgBot,      // Store the initialized TgBot
		// FbBot:      fbBot,      // Store the initialized FbBot
		// GeneralBot: generalBot, // Store the GeneralBot for /api/message route
	}
}
