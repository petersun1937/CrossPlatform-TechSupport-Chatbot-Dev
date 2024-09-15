package handlers

import (
	"Tg_chatbot/bot"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleGeneralWebhook handles incoming POST requests from the frontend
func HandleGeneralWebhook(c *gin.Context, b bot.GeneralBot) {
	// Parse the incoming request from the frontend and extract the message
	/*var req struct {
		Message   string `json:"message"`
		SessionID string `json:"sessionID"`
	}

	// Bind the incoming request body to the req struct
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("failed to bind request: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request"})
		return
	}*/

	// Delegate the handling of the message to the generalBot
	b.HandleGeneralMessage(c)

	// Return an OK status
	c.Status(http.StatusOK)
}
