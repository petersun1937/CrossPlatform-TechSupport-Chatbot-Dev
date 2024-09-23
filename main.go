package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Tg_chatbot/bot"
	config "Tg_chatbot/configs"
	"Tg_chatbot/database"
	"Tg_chatbot/server"
	"Tg_chatbot/service"
)

/*
main -> line (handle_line_msg)  /ask_line (questions)         -> /ask (questions) handle_internal_msg -> process command -> response msg -> response line or response tg
main -> telegram (handle_telegram_msg) /ask_tg (questions)

*/

func createBots(conf *config.Config, srv *service.Service) map[string]bot.Bot {
	// Initialize bots
	lineBot, err := bot.NewLineBot(conf, srv)
	if err != nil {
		//log.Fatal("Failed to initialize LINE bot:", err)
		fmt.Printf("Failed to initialize LINE bot: %s", err.Error())
	}

	tgBot, err := bot.NewTGBot(conf, srv)
	if err != nil {
		//log.Fatal("Failed to initialize Telegram bot:", err)
		fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	}

	fbBot, err := bot.NewFBBot(conf, srv)
	if err != nil {
		log.Fatalf("Failed to create Facebook bot: %v", err)
	}

	return map[string]bot.Bot{
		"line": lineBot,
		"tg":   tgBot,
		"fb":   fbBot,
	}
}

func main() {

	// Initialize config (only once)
	conf := config.GetConfig()

	// Initialize database
	db := database.NewDatabase(conf)
	if err := db.Init(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	// Initialize service
	srv := service.NewService(db)

	// initialize bots
	bots := createBots(conf, srv)

	// Initialize the app (app acts as the central hub for the application, holds different initialized values)
	app := server.NewApp(conf, srv, bots)
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

	// Initialize bots
	// lineBot, err := bot.NewLineBot(conf, srv)
	// if err != nil {
	// 	//log.Fatal("Failed to initialize LINE bot:", err)
	// 	fmt.Printf("Failed to initialize LINE bot: %s", err.Error())
	// }

	// tgBot, err := bot.NewTGBot(conf, srv)
	// if err != nil {
	// 	//log.Fatal("Failed to initialize Telegram bot:", err)
	// 	fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	// }

	// fbBot, err := bot.NewFBBot(conf, srv)
	// if err != nil {
	// 	log.Fatalf("Failed to create Facebook bot: %v", err)
	// }

	// Set webhook for Telegram using the ngrok URL (The set webhook step for LINE is done on their platform)
	// if err := tgBot.SetWebhook(conf.TelegramWebhookURL); err != nil {
	// 	log.Fatal("Failed to set Telegram webhook:", err)
	// }

	/*op bots := []bot.Bot{
		lineBot,
		tgBot,
		fbBot,
	}*/

	// initialize database
	// dbstr := os.Getenv("DATABASE_URL")
	// database.InitPostgresDB(dbstr) // Initialize the database connection (defined in package "DB")
	/*webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable not set")
	}*/

	// initialize http server routes from app struct
	// go app.RunRoutes(conf, srv)

	// // running bots
	// for _, bot := range bots {
	// 	if err := bot.Run(); err != nil {
	// 		//log.Fatal("running bot failed:", err)
	// 		fmt.Printf("running bot failed: %s", err.Error())
	// 	}
	// }

}
