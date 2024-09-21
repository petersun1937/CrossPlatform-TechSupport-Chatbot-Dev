package server

import (
	"Tg_chatbot/bot"
	config "Tg_chatbot/configs"
	"Tg_chatbot/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config     *config.Config
	Service    *service.Service // for database operations
	Router     *gin.Engine
	LineBot    bot.LineBot
	TgBot      bot.TgBot
	FbBot      bot.FbBot
	GeneralBot bot.GeneralBot // custom frontend
}

func NewApp(conf *config.Config, srv *service.Service) *App {
	// Initialize the line bot
	lineBot, err := bot.NewLineBot(conf, srv)
	if err != nil {
		log.Fatal("Failed to initialize LINE bot:", err)
	}

	// Initialize the tg bot
	tgBot, err := bot.NewTGBot(conf, srv)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	}

	// Initialize the tg bot
	fbBot, err := bot.NewFBBot(conf, srv)
	if err != nil {
		fmt.Printf("Failed to initialize Messenger bot: %s", err.Error())
	}

	// Initialize the general bot
	generalBot := bot.NewGeneralBot(nil)

	return &App{
		Config:     conf,
		Service:    srv,
		Router:     gin.Default(),
		LineBot:    lineBot,    // Store the initialized LineBot
		TgBot:      tgBot,      // Store the initialized TgBot
		FbBot:      fbBot,      // Store the initialized FbBot
		GeneralBot: generalBot, // Store the GeneralBot for /api/message route
	}
}
