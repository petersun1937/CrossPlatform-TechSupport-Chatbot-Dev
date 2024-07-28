package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Telegram API URL
const telegramAPIURL = "https://api.telegram.org/bot"

// Send a response to a Telegram chat
func sendTelegramResponse(chatID int64, response string) {
	// Construct the URL for the Telegram API request
	url := telegramAPIURL + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	// Create the message payload
	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    response,
	}

	// Marshal the message payload to JSON
	jsonMessage, _ := json.Marshal(message)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}
	defer resp.Body.Close()

	// Log the response
	log.Printf("Response sent to chat ID %d", chatID)
}

// Global variable to hold the bot instance
var Bot *tgbotapi.BotAPI

// Menu texts
var (
	FirstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	SecondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	NextButton     = "Next"
	BackButton     = "Back"
	TutorialButton = "Tutorial"

	// Keyboard layout for the first menu. One button, one row
	FirstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(NextButton, NextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	SecondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BackButton, BackButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(TutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

// Receive updates from the Telegram API and handle them
func ReceiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			HandleTelegramUpdate(update)
		}
	}
}

// Handle incoming Telegram updates
func HandleTelegramUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		handleMessage(update.Message)
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
	}
}

// Handle incoming messages
func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("%s wrote %s", user.FirstName, text)
	log.Printf("Chat ID: %d", message.Chat.ID) // Log the Chat ID

	var err error
	if strings.HasPrefix(text, "/") {
		// Handle commands
		err = handleCommand(message.Chat.ID, text)
	} else if screaming && len(text) > 0 {
		// If screaming mode is on, send the text in uppercase
		msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(text))
		msg.Entities = message.Entities
		_, err = Bot.Send(msg)
	} else {
		// Copy the message without the sender's name
		copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		_, err = Bot.CopyMessage(copyMsg)
	}

	if err != nil {
		log.Printf("An error occurred: %s", err.Error())
	}
}

// Global variable to store screaming state
var screaming bool

// Handle incoming commands
func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/start":
		msg := tgbotapi.NewMessage(chatId, "Welcome to the bot!")
		_, err = Bot.Send(msg)
	case "/scream":
		screaming = true // Enable screaming mode

	case "/whisper":
		screaming = false // Disable screaming mode

	case "/menu":
		err = sendMenu(chatId) // Send a menu to the chat

	case "/help":
		msg := tgbotapi.NewMessage(chatId, "Here are some commands you can use: /start, /help, /scream, /whisper, /menu")
		_, err = Bot.Send(msg)

	case "/custom":
		sendTelegramResponse(chatId, "This is a custom response!") // Send a custom response
	default:
		msg := tgbotapi.NewMessage(chatId, "I don't know that command")
		_, err = Bot.Send(msg)
	}

	return err
}

// Handle button clicks in inline keyboards
func handleButton(query *tgbotapi.CallbackQuery) {
	var text string
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if query.Data == NextButton {
		text = SecondMenu
		markup = SecondMenuMarkup
	} else if query.Data == BackButton {
		text = FirstMenu
		markup = FirstMenuMarkup
	}

	// Send a callback response
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	Bot.Send(callbackCfg)

	// Edit the message text and markup to reflect the button click
	msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	Bot.Send(msg)
}

// Send a menu to the chat
func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, FirstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = FirstMenuMarkup
	_, err := Bot.Send(msg)
	return err
}
