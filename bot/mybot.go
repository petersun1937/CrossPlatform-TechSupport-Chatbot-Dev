package bot

import (
	"Tg_chatbot/service"
	"context"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
)

type GeneralBot interface {
	Run() error
	HandleGeneralMessage(context *gin.Context)
	//SetWebhook(webhookURL string) error
}

type generalBot struct {
	// Add any common fields if necessary, like configuration
	*BaseBot
	ctx context.Context
	//config map[string]string
}

func NewGeneralBot(service *service.Service) *generalBot {
	baseBot := &BaseBot{
		Platform: GENERAL,
		Service:  service,
	}

	return &generalBot{
		BaseBot: baseBot,
		ctx:     context.Background(),
	}
}
func (b *generalBot) Run() error {
	// Implement logic for running the bot
	fmt.Println("General bot is running...")
	return nil
}

func (b *generalBot) HandleGeneralMessage(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionID"`
		Message   string `json:"message"`
	}

	// Parse the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Invalid request: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Store the context (to use later for sending the response)
	storeContext(req.SessionID, c)

	/*user := message.From
	if user == nil {
		return
	}

	userIDStr := strconv.FormatInt(user.ID, 10)
	fmt.Printf("User ID: %s \n", userIDStr)

	token, err := b.validateAndGenerateToken(userIDStr, user)
	if err != nil {
		fmt.Printf("Error validating user: %s", err.Error())
		return
	}

	if token != nil {
		b.sendTelegramMessage(message.Chat.ID, "Welcome! Your access token is: "+*token)
	} else {
		b.processUserMessage(message, user.FirstName, message.Text)
	}*/

	// Process and send the message
	b.ProcessUserMessage(req.SessionID, req.Message)

	// Send the response back to the frontend using sendResponse
	/*err = b.sendFrontendMessage(c, response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
		return
	}*/

	/*if err := b.sendResponse(req.SessionID, response); err != nil {
		fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
		return
	}*/
}

// ProcessUserMessage processes incoming messages
func (b *generalBot) ProcessUserMessage(sessionID string, message string) {
	var response string
	var err error

	fmt.Printf("Received message %s \n", message)
	fmt.Printf("Chat ID: %s \n", sessionID)

	// Example: Handle commands
	if strings.HasPrefix(message, "/") {
		response, err = handleCommand("", message, b) // TODO: Add ChatID instead of empty string
		if err != nil {
			fmt.Printf("An error occurred: %s \n", err.Error())
			response = "An error occurred while processing your command."
		}
	} else if screaming && len(message) > 0 {
		response = strings.ToUpper(message)
	} else {
		//response = fmt.Sprintf("You said: %s", message)
		handleMessageDialogflow(GENERAL, sessionID, message, b)
	}

	if response != "" {
		err := b.sendResponse(sessionID, response)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
			fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
		}
	}

}

func (b *generalBot) sendResponse(identifier interface{}, response string) error {
	// Perform type assertion to convert identifier to string
	if sessionID, ok := identifier.(string); ok {
		// Retrieve context using the sessionID
		c, err := getContext(sessionID)
		if err != nil {
			return fmt.Errorf("failed to retrieve context for sessionID: %s, error: %w", sessionID, err)
		}
		// Call sendFrontendMessage using the retrieved context
		return b.sendFrontendMessage(c, response)
	}
	return fmt.Errorf("invalid identifier type, expected string")
}

var sessionContextMap = make(map[string]*gin.Context)

// Store the context when the session starts
func storeContext(sessionID string, c *gin.Context) {
	sessionContextMap[sessionID] = c
}

// Retrieve the context using sessionID when you need to send a response
func getContext(sessionID string) (*gin.Context, error) {
	if context, ok := sessionContextMap[sessionID]; ok {
		return context, nil
	}
	return nil, fmt.Errorf("no context found for session ID %s", sessionID)
}

func (b *generalBot) sendFrontendMessage(c *gin.Context, message string) error {
	c.JSON(http.StatusOK, gin.H{
		"response": message,
	})
	return nil
}

func (b *generalBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {
	// Send the response to the respective platform or frontend
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if text := msg.GetText(); text != nil {
			return b.sendResponse(identifier, text.Text[0])
		}
	}
	return fmt.Errorf("invalid identifier for frontend or platform")
}

func (b *generalBot) sendMenu(identifier interface{}) error {
	if sessionID, ok := identifier.(string); ok {
		// Logic to send menu to the frontend user
		// For example, return a message to the frontend via the API response
		fmt.Printf("Sending menu to frontend user with session ID: %s\n", sessionID)
		return nil
	}
	return fmt.Errorf("invalid identifier type for frontend platform")
}
