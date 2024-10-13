package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/handlers"
	"crossplatform_chatbot/server"
	"crossplatform_chatbot/service"
)

/*
	Architecture
	[Prerequisites]
	- config

	[Application]
	- Server (handler/gin + service/database/bots...)

	[API]
	- Handler
	- Gin Server

	[Core Logic]
	- logic (service)
	- database (postgres)
	- bots
*/

func main() {

	// Initialize config (only once)
	conf := config.GetConfig()

	// Initialize database
	db := database.NewDatabase(conf)
	if err := db.Init(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	// Initialize service
	srv := service.NewService(conf.BotConfig, db)

	// Initialize handler
	handler := handlers.NewHandler(srv)

	// Initialize the app (app acts as the central hub for the application, holds different initialized values)
	app := server.NewApp(*conf, handler)
	if err := app.Run(); err != nil {
		log.Fatal("Failed to run the app:", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1) // creates a channel named quit

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // tells the program to listen for specific signals (SIGINT and SIGTERM) and send them to the quit channel.
	<-quit                                               // channel receive operation; blocking/waiting until a signal is received in quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // create context with timeout
	defer cancel()                                                          // ensure the context is canceled when the function exists

	if err := server.Shutdown(ctx); err != nil { // graceful shutdown
		log.Fatal("Server Shutdown: ", err)
	}

	fmt.Println("Server exiting")

}

// func createBots(conf *config.Config, srv *service.Service) map[string]bot.Bot {
// 	// Initialize bots
// 	lineBot, err := bot.NewLineBot(conf, srv)
// 	if err != nil {
// 		//log.Fatal("Failed to initialize LINE bot:", err)
// 		fmt.Printf("Failed to initialize LINE bot: %s", err.Error())
// 	}

// 	tgBot, err := bot.NewTGBot(conf, srv)
// 	if err != nil {
// 		//log.Fatal("Failed to initialize Telegram bot:", err)
// 		fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
// 	}

// 	fbBot, err := bot.NewFBBot(conf, srv)
// 	if err != nil {
// 		log.Fatalf("Failed to create Facebook bot: %v", err)
// 	}

// 	igBot, err := bot.NewIGBot(conf, srv)
// 	if err != nil {
// 		log.Fatalf("Failed to create Instagram bot: %v", err)
// 	}

// 	generalBot, err := bot.NewGeneralBot(conf, srv)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize General bot: %v", err)
// 	}

// 	return map[string]bot.Bot{
// 		"line":    lineBot,
// 		"tg":      tgBot,
// 		"fb":      fbBot,
// 		"ig":      igBot,
// 		"general": generalBot,
// 	}
// }
