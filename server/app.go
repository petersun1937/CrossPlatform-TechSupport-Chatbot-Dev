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
	Config  *config.Config
	Service *service.Service
	Router  *gin.Engine
	LineBot bot.LineBot
	TgBot   bot.TgBot
}

func NewApp(conf *config.Config, srv *service.Service) *App {
	lineBot, err := bot.NewLineBot(conf, srv)
	if err != nil {
		log.Fatal("Failed to initialize LINE bot:", err)
	}

	tgBot, err := bot.NewTGBot(conf, srv)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	}

	return &App{
		Config:  conf,
		Service: srv,
		Router:  gin.Default(),
		LineBot: lineBot, // Store the initialized LineBot
		TgBot:   tgBot,   // Store the initialized TgBot
	}
}
