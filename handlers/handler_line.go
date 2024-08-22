package handlers

import (
	"Tg_chatbot/database"
	"Tg_chatbot/models"
	"Tg_chatbot/utils"
	"Tg_chatbot/utils/token"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

// HandleLineWebhook handles incoming POST requests from the Line platform
func HandleLineWebhook(c *gin.Context) {
	// Parse the incoming request from the Line platform and extract the events
	events, err := utils.LineBot.ParseRequest(c.Request)
	if err != nil {
		// If the request has an invalid signature, return a 400 Bad Request error
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		// If there is any other error during parsing, return a 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse request"})
		return
	}

	// Loop through each event received from the Line platform
	for _, event := range events {
		// Check if the event is a message event
		if event.Type == linebot.EventTypeMessage {
			// Switch on the type of message (could be text, image, video, audio etc., only support text now)
			switch message := event.Message.(type) {
			// If the message is a text message, process it using handleLineMessage
			case *linebot.TextMessage:
				handleLineMessage(event, message)
			}
		}
	}
	c.Status(http.StatusOK)
}

// Processes incoming updates from Linebot
/*func HandleLineUpdate(event *linebot.Event) {
	if event.Type == linebot.EventTypeMessage {
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			handleLineMessage(event, message)
		}
	}
}*/

// Process incoming messages from users
func handleLineMessage(event *linebot.Event, message *linebot.TextMessage) {

	// Get user profile information from Line
	userProfile, err := utils.LineBot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		fmt.Printf("Error fetching user profile: %v\n", err)
		return
	}

	userID := event.Source.UserID
	text := message.Text

	fmt.Printf("User ID: %s \n", userID)

	//userIDInt, _ := strconv.ParseInt(userID, 10, 64)

	// Log the received message for debugging
	fmt.Printf("Received message: %s \n", text)

	// Check if the user exists in the database
	var dbUser models.User
	err = database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error

	// If the user does not exist, create a new user record
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dbUser = models.User{
				UserID:       userID, // LINE UserID as unique identifier
				UserName:     userProfile.DisplayName,
				FirstName:    "", // LINE doesn't provide firstname and lastname
				LastName:     "", // LINE doesn't provide firstname and lastname
				LanguageCode: userProfile.Language,
			}
			err = database.DB.Create(&dbUser).Error
			if err != nil {
				fmt.Printf("Error creating user: %s", err.Error())
				return
			}

			// Generate a JWT token for the new user
			token, err := token.GenerateToken(userID, "user") // Convert userID to int if needed
			if err != nil {
				fmt.Printf("Error generating JWT: %s", err.Error())
				return
			}

			// Send the token to the user
			msg := linebot.NewTextMessage("Welcome! Your access token is: " + token)
			if _, err := utils.LineBot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
				fmt.Printf("Error sending token message: %s \n", err.Error())
			}
		} else {
			// Handle other types of errors
			fmt.Printf("Error retrieving user: %s", err.Error())
		}
	}
	if strings.HasPrefix(text, "/") {
		err := handleLineCommand(event, text)
		if err != nil {
			fmt.Printf("An error occurred: %s \n", err.Error())
		}
		return
	}
	// If not a command, process the message using processLineMessage function
	/*response := processMessage(text)
	if _, err := utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
		fmt.Printf("An error occurred: %s \n", err.Error())
	}*/

	// Process the message using Dialogflow
	handleLineMessageDF(event, message)
}

func handleLineMessageDF(event *linebot.Event, message *linebot.TextMessage) {
	projectID := "testagent-mkyg"    // dialogflow project id
	sessionID := event.Source.UserID // Use user ID as session ID
	languageCode := "en"

	response, err := utils.DetectIntentText(projectID, sessionID, message.Text, languageCode)
	if err != nil {
		fmt.Printf("Error detecting intent: %v\n", err)
		return
	}

	handleDialogflowResponse(response, LINE, event)
}

func handleLineCommand(event *linebot.Event, command string) error {
	var err error

	switch command {
	case "/start":
		if _, err = utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Welcome to the bot!")).Do(); err != nil {
			fmt.Print(err)
		}
	case "/scream":
		utils.Screaming = true // Enable screaming mode
		if _, err = utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Scream mode enabled!")).Do(); err != nil {
			fmt.Print(err)
		}
	case "/whisper":
		utils.Screaming = false // Disable screaming mode
		if _, err = utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Scream mode disabled!")).Do(); err != nil {
			fmt.Print(err)
		}
	case "/menu":
		err = utils.SendLineMenu(event.ReplyToken) // Send a menu to the chat
	case "/help":
		if _, err = utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Here are some commands you can use: /start, /help, /scream, /whisper, /menu")).Do(); err != nil {
			fmt.Print(err)
		}
	default:
		if _, err = utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("I don't know that command")).Do(); err != nil {
			fmt.Print(err)
		}
	}

	return err
}
