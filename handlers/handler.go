package handlers

import (
	"Tg_chatbot/models"
	"Tg_chatbot/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handles incoming updates from the Telegram webhook
func HandleTelegramWebhook(c *gin.Context) {
	var update models.TelegramUpdate
	// Bind the incoming JSON payload to the update model
	if err := c.BindJSON(&update); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Printf("Received message from Telegram: %s", update.Message.Text)
	// Convert the custom update model to the tgbotapi.Update type and handle it
	utils.HandleTelegramUpdate(tgbotapi.Update{
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

/*func HandleTelegramWebhook(c *gin.Context) {
	var update models.TelegramUpdate
	if err := c.BindJSON(&update); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Printf("Received message from Telegram: %s", update.Message.Text)
	response := processMessage(update.Message.Text)
	sendTelegramResponse(update.Message.Chat.ID, response)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}*/

// HandleCustomMessage processes incoming messages from the custom frontend
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
}
