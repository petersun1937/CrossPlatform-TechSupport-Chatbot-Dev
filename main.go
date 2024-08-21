package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"Tg_chatbot/database"
	"Tg_chatbot/handlers"
	"Tg_chatbot/server"
	"Tg_chatbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load("configs/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	// Get bot token and webhookurl from environment variable
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	dbstr := os.Getenv("DATABASE_URL")
	database.InitPostgresDB(dbstr) // Initialize the database connection (defined in package "DB")
	/*webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("WEBHOOK_URL environment variable not set")
	}*/

	// Initialize Telegram Bot
	utils.TgBot, err = tgbotapi.NewBotAPI(botToken) // creates new BotAPI instance using the token
	// utils.Bot: global variable (defined in the utils package) that holds the reference to the bot instance.

	if err != nil {
		log.Panic(err)
	}
	utils.TgBot.Debug = false

	// Initialize Linebot
	err = utils.InitLineBot(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal("Error initializing Linebot:", err)
	}

	// Set webhook URL
	//_, err = utils.Bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))

	// Set webhook with certificate
	/*certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")
	webhook, err := tgbotapi.NewWebhookWithCert(webhookURL+"/"+utils.Bot.Token, tgbotapi.FilePath(certFile))
	if err != nil {
		log.Fatal(err)
	}

	_, err = utils.Bot.Request(webhook)
	if err != nil {
		log.Panic(err)
	}

	info, err := utils.Bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	// Start the server to listen for webhook updates
	go http.ListenAndServeTLS("0.0.0.0:8443", certFile, keyFile, nil)*/

	// Initialize and start the server
	// using a go routine allows the program to handle multiple tasks simultaneously without blocking.
	go server.RunRoutes()

	// Create a new cancellable background context
	// This provides a context (ctx) that can be passed around to different functions or goroutines.
	// context.Background() is often used as the root context for new goroutines when no specific request or context is available.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Create a new update configuration (long polling)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates channel
	updates := utils.TgBot.GetUpdatesChan(u)

	// Start receiving updates (go routine)
	go handlers.ReceiveUpdates(ctx, updates)

	// Wait for a newline symbol, then cancel handling updates (for this to work, run with cmd)
	fmt.Println("Bot is running. Press Enter to stop.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

/*
   func main() {

   	// Load environment variables from .env file

   	if err := godotenv.Load("configs/.env"); err != nil {
   		panic("Error loading .env file")
   	}

   	//dbstr := os.Getenv("DATABASE_URL")
   	//database.InitPostgresDB(dbstr) // Initialize the database connection (defined in package "DB")

   	server.RunRoutes()

   }*/
