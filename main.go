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

	"github.com/joho/godotenv"
)

/*
main -> line (handle_line_msg)  /ask_line (questions)         -> /ask (questions) handle_internal_msg -> process command -> response msg -> response line or response tg
main -> telegram (handle_telegram_msg) /ask_tg (questions)

*/

func main() {
	// Load environment variables
	err := godotenv.Load("configs/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// load configs
	conf := config.NewConfig(dbURL)

	// database
	db := database.NewDatabase2(conf)
	if err := db.Init(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	// service
	srv := service.NewService(db)

	// Initialize the app
	app := server.NewApp(conf, srv)

	// server
	// server := NewServer()
	// server.Run()

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
	bots := []bot.Bot{
		lineBot,
		tgBot,
		//bot.NewLineBot(conf),
		//bot.NewTGBot(conf, srv),
	}

	// initialize database
	// dbstr := os.Getenv("DATABASE_URL")
	// database.InitPostgresDB(dbstr) // Initialize the database connection (defined in package "DB")
	/*webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable not set")
	}*/

	// initialize http server
	go app.RunRoutes(conf, srv)

	// running bots
	for _, bot := range bots {
		if err := bot.Run(); err != nil {
			//log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
		}
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	fmt.Println("Server exiting")

}

// func main2() {
// 	// Load environment variables
// 	err := godotenv.Load("configs/.env")
// 	if err != nil {
// 		panic("Error loading .env file")
// 	}

// 	// Get bot token and webhookurl from environment variable
// 	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
// 	if botToken == "" {
// 		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
// 	}

// 	dbstr := os.Getenv("DATABASE_URL")
// 	database.InitPostgresDB(dbstr) // Initialize the database connection (defined in package "DB")
// 	/*webhookURL := os.Getenv("WEBHOOK_URL")
// 	if webhookURL == "" {
// 		log.Fatal("WEBHOOK_URL environment variable not set")
// 	}*/

// 	Initialize Telegram Bot
// 	utils.TgBot, err = tgbotapi.NewBotAPI(botToken) // create new BotAPI instance using the token
// 	// utils.TgBot: global variable (defined in the utils package) that holds the reference to the bot instance.
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	// Initialize Linebot
// 	utils.LineBot, err = linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN")) // create new BotAPI instance using the channel token and secret
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	Set webhook URL
// 	_, err = utils.Bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))

// 	Set webhook with certificate
// 	certFile := os.Getenv("CERT_FILE")
// 	keyFile := os.Getenv("KEY_FILE")
// 	webhook, err := tgbotapi.NewWebhookWithCert(webhookURL+"/"+utils.Bot.Token, tgbotapi.FilePath(certFile))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = utils.Bot.Request(webhook)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	info, err := utils.Bot.GetWebhookInfo()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if info.LastErrorDate != 0 {
// 		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
// 	}

// 	// Start the server to listen for webhook updates
// 	go http.ListenAndServeTLS("0.0.0.0:8443", certFile, keyFile, nil)

// 	// Initialize and start the server
// 	using a go routine allows the program to handle multiple tasks simultaneously without blocking.
// 	go server.RunRoutes()

// 	Create a new cancellable background context
// 	This provides a context (ctx) that can be passed around to different functions or goroutines.
// 	context.Background() is often used as the root context for new goroutines when no specific request or context is available.
// 	ctx := context.Background()
// 	deadline := time.Now().Add(time.Second * 30)
// 	ctx, cancel := context.WithDeadline(ctx, deadline)
// 	ddd, ok := ctx.Deadline()
// 	cancel()
// 	ctx, cancel := context.WithCancel(ctx)

// 	Create a new update configuration with offset of 0
// 	Using 0 means it will start fetching updates from the beginning.
// 	u := tgbotapi.NewUpdate(0)
// 	u.Timeout = 60 // timeout for long polling set to 60 s

// 	// Get updates channel to start long polling to receive updates.
// 	// The channel will be continuously fed with new Update objects from Telegram.
// 	updates := utils.TgBot.GetUpdatesChan(u)

// 	// Use go routine to continuously process received updates from the updates channel
// 	go handlers.ReceiveUpdates(ctx, updates)

// 	Wait for a newline symbol, then cancel handling updates (for cancel to work, run with cmd)
// 	fmt.Println("Bot is running. Press Enter to stop.")
// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
// 	cancel()
// }
