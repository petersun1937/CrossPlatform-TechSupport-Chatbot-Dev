package handlers

import (
	"Tg_chatbot/models"
	"Tg_chatbot/utils"
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ReceiveUpdates receives updates from Telegram API and handles them
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

// HandleTelegramUpdate processes incoming updates from Telegram
func HandleTelegramUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		handleButton(update.CallbackQuery)
	}
}

// handleMessage processes incoming messages from users
func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("Received message from %s: %s", user.FirstName, text)
	log.Printf("Chat ID: %d", message.Chat.ID) // Log the Chat ID

	var err error
	if strings.HasPrefix(text, "/") {
		// Handle commands
		err = handleCommand(message.Chat.ID, text)
	} else if utils.Screaming && len(text) > 0 {
		// If screaming mode is on, send the text in uppercase
		msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(text))
		msg.Entities = message.Entities
		_, err = utils.Bot.Send(msg)
	} else {
		// Process the message using processMessage function
		response := processMessage(text)
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		_, err = utils.Bot.Send(msg)

		// Copy the message without the sender's name
		//msg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		//_, err = utils.Bot.CopyMessage(msg)
	}

	if err != nil {
		log.Printf("An error occurred: %s", err.Error())
	}
}

func processMessage(message string) string {
	message = strings.ToLower(message)
	switch {
	case strings.Contains(message, "hello"):
		return "Hello! How may I help you?"
	case strings.Contains(message, "help"):
		return "Here are some commands you can use: /start, /help, /scream, /whisper, /menu"
	case strings.Contains(message, "scream"):
		utils.Screaming = true
		return "Scream mode enabled!"
	case strings.Contains(message, "whisper"):
		utils.Screaming = false
		return "Scream mode disabled!"
	default:
		return "I'm sorry, I didn't understand that. Type /help to see what I can do."
	}
}

// handleCommand processes commands sent by users
func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/start":
		msg := tgbotapi.NewMessage(chatId, "Welcome to the bot!")
		_, err = utils.Bot.Send(msg)
	case "/scream":
		utils.Screaming = true // Enable screaming mode
		msg := tgbotapi.NewMessage(chatId, "Scream mode enabled!")
		_, err = utils.Bot.Send(msg)

	case "/whisper":
		utils.Screaming = false // Disable screaming mode
		msg := tgbotapi.NewMessage(chatId, "Scream mode disabled!")
		_, err = utils.Bot.Send(msg)

	case "/menu":
		err = utils.SendMenu(chatId) // Send a menu to the chat

	case "/help":
		msg := tgbotapi.NewMessage(chatId, "Here are some commands you can use: /start, /help, /scream, /whisper, /menu")
		_, err = utils.Bot.Send(msg)

	case "/custom":
		utils.SendTelegramResponse(chatId, "This is a custom response!") // Send a custom response
	default:
		msg := tgbotapi.NewMessage(chatId, "I don't know that command")
		_, err = utils.Bot.Send(msg)
	}

	return err
}

// handleButton processes callback queries from inline keyboard buttons
func handleButton(query *tgbotapi.CallbackQuery) {
	var text string
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if query.Data == utils.NextButton {
		text = utils.SecondMenu
		markup = utils.SecondMenuMarkup
	} else if query.Data == utils.BackButton {
		text = utils.FirstMenu
		markup = utils.FirstMenuMarkup
	}

	// Send a callback response
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	utils.Bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	utils.Bot.Send(msg)
}

// handles incoming updates from the Telegram webhook
func HandleTelegramWebhook(c *gin.Context) {
	var update models.TelegramUpdate
	if err := c.BindJSON(&update); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	HandleTelegramUpdate(tgbotapi.Update{
		UpdateID: update.UpdateID,
		Message: &tgbotapi.Message{
			MessageID: update.Message.MessageID,
			From: &tgbotapi.User{
				ID:           update.Message.From.ID,
				IsBot:        update.Message.From.IsBot,
				FirstName:    update.Message.From.FirstName,
				LastName:     update.Message.From.LastName,
				UserName:     update.Message.From.UserName,
				LanguageCode: update.Message.From.LanguageCode,
			},
			Chat: &tgbotapi.Chat{
				ID:        update.Message.Chat.ID,
				Type:      update.Message.Chat.Type,
				Title:     update.Message.Chat.Title,
				UserName:  update.Message.Chat.UserName,
				FirstName: update.Message.Chat.FirstName,
				LastName:  update.Message.Chat.LastName,
			},
			Date: update.Message.Date,
			Text: update.Message.Text,
		},
	})
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// processes custom messages
/*
func HandleCustomMessage(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}
	// Bind the incoming JSON payload to the request struct
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Printf("Received message from frontend: %s", req.Message)
	// Process the message and generate a response
	response := processMessage(req.Message)
	c.JSON(http.StatusOK, gin.H{"response": response})
}

// generates a response for the given message
func processMessage(message string) string {
	// Placeholder for NLP model integration
	return "Processed message: " + message
}*/
