package server

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/handlers"
)

type App struct { //TODO: app or GetConfig for config?
	Config config.Config
	// Service *service.Service // for database operations
	// Router  *gin.Engine
	// Bots map[string]bot.Bot
	//LineBot    bot.LineBot
	//TgBot      bot.TgBot
	//FbBot      bot.FbBot
	//GeneralBot bot.GeneralBot // custom frontend

	Server *Server
}

func (a App) Run() error {
	// initialize http server routes from app struct
	a.Server.V2Start()
	// go a.RunRoutes(a.Config, a.Service)

	// // running bots
	// for _, bot := range a.Bots {
	// 	if err := bot.Run(); err != nil {
	// 		// log.Fatal("running bot failed:", err)
	// 		fmt.Printf("running bot failed: %s", err.Error())
	// 		return err
	// 	}
	// }

	return nil
}

func NewApp(conf config.Config, handler *handlers.Handler) *App {
	return &App{
		Config: conf,
		// Service: srv,
		// Router:  gin.Default(),
		// Bots: bots,
		// LineBot:    lineBot,    // Store the initialized LineBot
		// TgBot:      tgBot,      // Store the initialized TgBot
		// FbBot:      fbBot,      // Store the initialized FbBot
		// GeneralBot: generalBot, // Store the GeneralBot for /api/message route
		Server: New(conf.ServerConfig, handler),
	}
}
