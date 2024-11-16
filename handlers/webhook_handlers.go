package handlers

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

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

// const (
// 	LINE bot.Platform = iota
// 	TELEGRAM
// 	FACEBOOK
// 	INSTAGRAM
// 	GENERAL
// )

// func parsePlatformAndIdentifier(sessionID string) (bot.Platform, interface{}, error) {
// 	// Parse the sessionID to determine the platform and identifier
// 	// Example: you might encode sessionID as "LINE:userID" or "TELEGRAM:chatID"
// 	parts := strings.Split(sessionID, ":")
// 	if len(parts) != 2 {
// 		return 0, nil, fmt.Errorf("invalid sessionID format")
// 	}

// 	platformStr, identifier := parts[0], parts[1]
// 	switch platformStr {
// 	case "LINE":
// 		return LINE, identifier, nil
// 	case "TELEGRAM":
// 		return TELEGRAM, identifier, nil
// 	case "FACEBOOK":
// 		return FACEBOOK, identifier, nil
// 	case "GENERAL":
// 		return GENERAL, identifier, nil
// 	default:
// 		return 0, nil, fmt.Errorf("unknown platform")
// 	}
// }

// func (h *Handler) HandleDialogflowWebhook(c *gin.Context) {
// 	// Read the raw JSON body
// 	var rawRequestBody map[string]interface{}
// 	if err := c.BindJSON(&rawRequestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Marshal the raw JSON into bytes
// 	requestBytes, err := json.Marshal(rawRequestBody)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
// 		return
// 	}

// 	// Unmarshal the JSON bytes into the protobuf struct
// 	var request dialogflowpb.WebhookRequest
// 	if err := protojson.Unmarshal(requestBytes, &request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unmarshal to WebhookRequest"})
// 		return
// 	}

// 	// Immediate response to avoid Dialogflow timeouts
// 	ackResponse := gin.H{
// 		"fulfillmentText": "Processing your request. Please wait...",
// 	}
// 	c.JSON(http.StatusOK, ackResponse)

// 	// Process the request asynchronously to prevent timeouts
// 	go func() {
// 		response, err := h.Service.ProcessIntentWithRAG(&request)
// 		if err != nil {
// 			fmt.Printf("Error processing intent: %v\n", err)
// 			return
// 		}

// 		fmt.Printf("Processed response: %+v\n", response)
// 	}()
// }

// func (h *Handler) HandleDialogflowWebhook(c *gin.Context) {
// 	var request dialogflowpb.WebhookRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(400, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Process the webhook request here
// 	response, err := h.Service.ProcessIntentWithRAG(request)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Send back a response to Dialogflow
// 	c.JSON(200, response)
// }

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

// // HandleTelegramWebhook handles POST requests from Telegram.
//
//	func (h *Handler) HandleTelegramWebhook(c *gin.Context) {
//		var update tgbotapi.Update
//		if err := c.ShouldBindJSON(&update); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request"})
//			return
//		}
//		if err := h.Service.HandleTelegram(update); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.Status(http.StatusOK)
//	}
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
