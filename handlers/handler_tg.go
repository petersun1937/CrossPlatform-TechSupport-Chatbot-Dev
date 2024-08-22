package handlers

import (
	"Tg_chatbot/database"
	"Tg_chatbot/models"
	"Tg_chatbot/utils"
	"Tg_chatbot/utils/token"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// ReceiveUpdates receives updates from Telegram API and handles them
func ReceiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// "updates" is a channel that receives updates from the Telegram bot (e.g., messages, button clicks).
	// The bot's API sends these updates to the application, and the function processes them by handling the updates.

	for { // continuous loop to check for updates
		select {
		case <-ctx.Done(): // If context has been cancelled
			fmt.Println("Goroutine: Received cancel signal, stopping...")
			// exit the loop and stop the go routine
			return
		case update := <-updates: // Process incoming updates from Telegram
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

// Processes incoming messages from users
func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Convert user.ID from int64 to string
	userIDStr := strconv.FormatInt(user.ID, 10)

	fmt.Printf("User ID: %s \n", userIDStr)

	// Check if the user exists in the database
	var dbUser models.User
	err := database.DB.Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error

	// If the user does not exist, create a new user record
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist, create a new user record
			dbUser = models.User{
				UserID:       userIDStr,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				UserName:     user.UserName,
				LanguageCode: user.LanguageCode,
			}
			err = database.DB.Create(&dbUser).Error
			if err != nil {
				fmt.Printf("Error creating user: %s", err.Error())
				return
			}

			// Generate a JWT token for the new user
			token, err := token.GenerateToken(userIDStr, "user") // Ensure GenerateToken accepts string
			if err != nil {
				fmt.Printf("Error generating JWT: %s", err.Error())
				return
			}

			// Send the token to the user
			msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Your access token is: "+token)
			utils.TgBot.Send(msg)
		} else {
			// Handle other types of errors
			fmt.Printf("Error retrieving user: %s", err.Error())
		}

	} else {
		fmt.Printf("Received message from %s: %s \n", user.FirstName, text)
		fmt.Printf("Chat ID: %d \n", message.Chat.ID)

		var response string
		if strings.HasPrefix(text, "/") {
			// Handle commands
			//response, err = handleCommand(message.Chat.ID, text)
			response, err = handleCommand(TELEGRAM, message.Chat.ID, text)
			if err != nil {
				fmt.Printf("An error occurred: %s \n", err.Error())
				response = "An error occurred while processing your command."
			}
		} else if utils.Screaming && len(text) > 0 {
			// If screaming mode is on, send the text in uppercase
			response = strings.ToUpper(text)
		} else {
			// Process the message using processMessage function
			//response = processMessage(text)

			// Send the message to Dialogflow for processing
			handleTGMessageDialogflow(message)
			return
		}

		fmt.Printf("Response: '%s'\n", response)

		// Send the response if it's not empty
		if response != "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, response)
			_, err = utils.TgBot.Send(msg)
			if err != nil {
				fmt.Printf("An error occurred: %s \n", err.Error())
			}
		}
	}
}

// Handle messages with Dialogflow
func handleTGMessageDialogflow(message *tgbotapi.Message) {
	projectID := "testagent-mkyg"
	sessionID := strconv.FormatInt(message.Chat.ID, 10)
	languageCode := "en"

	// Send the user’s message to Dialogflow and receives a response.
	response, err := utils.DetectIntentText(projectID, sessionID, message.Text, languageCode)
	if err != nil {
		fmt.Printf("Error detecting intent: %v\n", err)
		return
	}

	// Process Dialogflow response and send it
	handleDialogflowResponse(response, TELEGRAM, message.Chat.ID)
}

/*

// processes incoming messages from users

	func handleMessage(message *tgbotapi.Message) {
		user := message.From
		text := message.Text

		if user == nil {
			return
		}

		fmt.Printf("Received message from %s: %s", user.FirstName, text)
		fmt.Printf("Chat ID: %d", message.Chat.ID) // Log the Chat ID

		var err error
		if strings.HasPrefix(text, "/") {
			// Handle commands
			err = handleCommand(message.Chat.ID, text)
		} else if utils.Screaming && len(text) > 0 {
			// If screaming mode is on, send the text in uppercase
			msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(text))
			msg.Entities = message.Entities
			_, err = utils.TgBot.Send(msg)
		} else {
			// Process the message using processMessage function
			response := processMessage(text)
			msg := tgbotapi.NewMessage(message.Chat.ID, response)
			_, err = utils.TgBot.Send(msg)

			// Copy the message without the sender's name
			//msg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
			//_, err = utils.TgBot.CopyMessage(msg)
		}

		if err != nil {
			fmt.Printf("An error occurred: %s", err.Error())
		}
	}
*/
/*
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
*/

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

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	utils.TgBot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	utils.TgBot.Send(msg)
}

/*
// handles incoming updates from the Telegram webhook
func HandleTelegramWebhook(c *gin.Context) {
	var update models.TelegramUpdate
	if err := c.BindJSON(&update); err != nil {
		fmt.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Extract user information from the update
	user := update.Message.From

	// Check if the user exists in the database
	var dbUser models.User
	err := database.DB.Where("user_id = ?", user.ID).First(&dbUser).Error

	if err != nil {
		// User doesn't exist, create a new user
		dbUser = models.User{
			UserID:       user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			UserName:     user.UserName,
			LanguageCode: user.LanguageCode,
		}
		err = database.DB.Create(&dbUser).Error
		if err != nil {
			fmt.Printf("Error creating user: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		// Generate a JWT token for the new user
		token, err := token.GenerateToken(int(user.ID), "user")
		if err != nil {
			fmt.Printf("Error generating JWT: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		// Send the token to the user via Telegram
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Your access token is: "+token)
		_, err = utils.TgBot.Send(msg)
		if err != nil {
			fmt.Printf("Error sending token to user: %s", err.Error())
		}
	}

	// Process the message
	HandleTelegramUpdate(tgbotapi.Update{
		UpdateID: update.UpdateID,
		Message: &tgbotapi.Message{
			MessageID: update.Message.MessageID,
			From:      update.Message.From,
			Chat:      update.Message.Chat,
			Date:      update.Message.Date,
			Text:      update.Message.Text,
		},
	})

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
*/

// processes custom messages

func HandleCustomMessage(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}
	// Bind the incoming JSON payload to the request struct
	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	fmt.Printf("Received message from frontend: %s", req.Message)
	// Process the message and generate a response
	response := processMessage(req.Message)
	c.JSON(http.StatusOK, gin.H{"response": response})
}

/*
// generates a response for the given message
func processMessage(message string) string {
	// Placeholder for NLP model integration
	return "Processed message: " + message
}*/
