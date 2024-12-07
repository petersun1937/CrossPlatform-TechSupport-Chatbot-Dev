package handlers

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

func (h *Handler) HandleLineWebhook(c *gin.Context) {
	if err := h.Service.HandleLine(c.Request); err != nil {
		// If the request has an invalid signature, return a 400 Bad Request error
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		// If there is any other error during parsing, return a 500 Internal Server Error
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse request"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) HandleTelegramWebhook(c *gin.Context) {
	var update tgbotapi.Update

	// Log raw request body
	body, _ := c.GetRawData()
	fmt.Println("Received Telegram update:", string(body))

	// Try to bind the JSON to the update struct
	if err := json.Unmarshal(body, &update); err != nil {
		fmt.Println("Failed to bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request"})
		return
	}

	// Process the update and log any error
	if err := h.Service.HandleTelegram(update); err != nil {
		fmt.Println("Error handling update:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Successfully processed update")
	c.Status(http.StatusOK)
}

// HandleMessengerWebhook handles POST requests from Facebook Messenger.
func (h *Handler) HandleMessengerWebhook(c *gin.Context) {
	var event bot.MessengerEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}
	if err := h.Service.HandleMessenger(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// HandleInstagramWebhook handles POST requests from Instagram.
func (h *Handler) HandleInstagramWebhook(c *gin.Context) {
	var event bot.InstagramEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
		return
	}
	if err := h.Service.HandleInstagram(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// VerifyInstagramWebhook verifies the webhook for Instagram Messaging (handles GET request)
func (h *Handler) VerifyInstagramWebhook(c *gin.Context) {
	// Load verification token from configuration or environment
	conf := config.GetConfig()
	verifyToken := conf.InstagramVerifyToken // Use Instagram-specific verify token

	// Check if the verify token matches
	if c.Query("hub.verify_token") == verifyToken {
		c.String(http.StatusOK, c.Query("hub.challenge"))
	} else {
		c.String(http.StatusForbidden, "Invalid verification token")
	}
}

// VerifyMessengerWebhook verifies the webhook for Facebook Messenger (handles GET request)
func (h *Handler) VerifyMessengerWebhook(c *gin.Context) {
	// Verify token from environment or configuration
	//verifyToken := os.Getenv("VERIFY_TOKEN")
	conf := config.GetConfig()
	verifyToken := conf.FacebookVerifyToken

	// Check if the verify token matches
	if c.Query("hub.verify_token") == verifyToken {
		c.String(http.StatusOK, c.Query("hub.challenge"))
	} else {
		c.String(http.StatusForbidden, "Invalid verification token")
	}
}

// HandlerGeneralBot handles incoming POST requests from the frontend
func (h *Handler) HandlerGeneralBot(c *gin.Context) {
	var req models.GeneralRequest

	// Parse the incoming request from the frontend and bind to the req struct
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("failed to bind request: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// // Store the context (to use later for sending the response) TODO to bot?
	b := h.Service.GetBot("general").(bot.GeneralBot)

	// Store the context using the sessionID
	b.StoreContext(req.SessionID, c)

	// Delegate the request to the service layer.
	if err := h.Service.HandleGeneral(req); err != nil {
		fmt.Printf("Error handling general request: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle request"})
		return
	}

	// Return an OK status after successfully processing the message
	c.Status(http.StatusOK)
}

// Webhook for Dialogflow (discarded, use direct response)
// func (h *Handler) HandleDialogflowWebhook(c *gin.Context) {
// 	// Read and unmarshal the request body into a protobuf struct
// 	var request dialogflowpb.WebhookRequest
// 	if err := c.BindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Immediate response to avoid Dialogflow timeouts
// 	ackResponse := gin.H{
// 		"fulfillmentText": "Processing your request. Please wait...",
// 	}
// 	c.JSON(http.StatusOK, ackResponse)

// 	// Process the intent asynchronously to prevent timeouts
// 	go func() {
// 		// Retrieve the session ID from Dialogflow request
// 		sessionID := request.Session
// 		platform, identifier, err := parsePlatformAndIdentifier(sessionID)
// 		if err != nil {
// 			fmt.Printf("Error parsing platform and identifier: %v\n", err)
// 			return
// 		}

// 		// Process the intent with RAG (retrieval-augmented generation)
// 		response, err := h.Service.ProcessIntentWithRAG(&request)
// 		if err != nil {
// 			fmt.Printf("Error processing intent: %v\n", err)
// 			return
// 		}

// 		// Send the response to the correct platform using the identifier
// 		err = h.sendFinalResponseToPlatform(platform, identifier, response.FulfillmentText)
// 		if err != nil {
// 			fmt.Printf("Error sending final response to platform: %v\n", err)
// 		}
// 	}()
// }