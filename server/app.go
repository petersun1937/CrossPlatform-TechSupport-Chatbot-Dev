package server

import (
	"Tg_chatbot/bot"
	config "Tg_chatbot/configs"
	"Tg_chatbot/service"
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config  *config.Config
	Service *service.Service
	Router  *gin.Engine
	LineBot bot.LineBot
}

func NewApp(conf *config.Config, srv *service.Service) *App {
	lineBot, err := bot.NewLineBot(conf, srv)
	if err != nil {
		log.Fatal("Failed to initialize LINE bot:", err)
	}

	return &App{
		Config:  conf,
		Service: srv,
		Router:  gin.Default(),
		LineBot: lineBot, // Store the initialized LineBot
	}
}
