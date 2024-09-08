package bot

import (
	config "Tg_chatbot/configs"
	"Tg_chatbot/models"
	"Tg_chatbot/service"
	"Tg_chatbot/utils/token"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"gorm.io/gorm"
)

type LineBot interface {
	Run() error
	ParseRequest(req *http.Request) ([]*linebot.Event, error)
	HandleLineMessage(event *linebot.Event, message *linebot.TextMessage)
}

type lineBot struct {
	*BaseBot
	secret     string
	token      string
	lineClient *linebot.Client
	service    *service.Service
}

func NewLineBot(conf *config.Config, service *service.Service) (*lineBot, error) {
	lineClient, err := linebot.New(conf.LineChannelSecret, conf.LineChannelToken)
	if err != nil {
		return nil, err
	}

	baseBot := &BaseBot{
		Platform: LINE,
		Service:  service,
	}

	return &lineBot{
		BaseBot:    baseBot,
		secret:     conf.LineChannelSecret,
		token:      conf.LineChannelToken,
		lineClient: lineClient,
		service:    service,
	}, nil
	/*return &lineBot{
		secret: conf.GetLineSecret(),
		token:  conf.GetLineToken(),
	}*/
}

func (b *lineBot) Run() error {
	// Initialize Linebot
	lineClient, err := linebot.New(b.secret, b.token) // create new BotAPI instance using the channel token and secret
	if err != nil {
		return err
	}

	b.lineClient = lineClient

	// Start the bot with webhook
	fmt.Println("Line bot is running with webhook!")

	return nil
}

func (b *lineBot) HandleLineMessage(event *linebot.Event, message *linebot.TextMessage) {

	// Retrieve and validate user profile
	userProfile, err := b.getUserProfile(event.Source.UserID)
	if err != nil {
		fmt.Printf("Error fetching user profile: %v\n", err)
		return
	}

	// Ensure user exists in the database
	userExists, err := b.ensureUserExists(userProfile, event, event.Source.UserID)
	if err != nil {
		fmt.Printf("Error ensuring user exists: %v\n", err)
		return
	}

	// If user didn't exist, a welcome message was sent, so return
	if !userExists {
		return
	}

	// Process the user's message
	response, err := b.processUserMessage(event, message.Text)
	if err != nil {
		fmt.Printf("Error processing user message: %v\n", err)
		return
	}

	// Send the response if it's not empty
	if response != "" {
		if err := b.sendLineMessage(event.ReplyToken, response); err != nil {
			fmt.Printf("Error sending response message: %v\n", err)
		}
	}
}

// Get user profile from Line
func (b *lineBot) getUserProfile(userID string) (*linebot.UserProfileResponse, error) {
	userProfile, err := b.lineClient.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

// Ensure the user exists in the database, create if not, return true if user existed
func (b *lineBot) ensureUserExists(userProfile *linebot.UserProfileResponse, event *linebot.Event, userID string) (bool, error) {
	var dbUser models.User
	//err := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
	err := b.service.GetDB().Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			dbUser = models.User{
				UserID:       userID,
				UserName:     userProfile.DisplayName,
				FirstName:    "", // LINE doesn't provide first and last names
				LastName:     "",
				LanguageCode: userProfile.Language,
			}
			//if err := database.DB.Create(&dbUser).Error; err != nil {
			if err := b.service.GetDB().Create(&dbUser).Error; err != nil {
				return false, fmt.Errorf("error creating user: %w", err)
			}

			// Generate a JWT token for the new user
			token, err := token.GenerateToken(userID, "user")
			if err != nil {
				return false, fmt.Errorf("error generating JWT: %w", err)
			}

			// Send welcome message with the token
			if err := b.sendLineMessage(event.ReplyToken, "Welcome! Your access token is: "+token); err != nil {
				return false, fmt.Errorf("error sending token message: %w", err)
			}
			return false, nil // User was just created and welcomed
		}
		return false, fmt.Errorf("error retrieving user: %w", err)
	}
	return true, nil // User already existed
}

// Process the user's message (commands or Dialogflow)
func (b *lineBot) processUserMessage(event *linebot.Event, text string) (string, error) {
	var response string
	var err error

	if strings.HasPrefix(text, "/") { // Command handling
		response, err = handleCommand(event, text, b)
		if err != nil {
			return "An error occurred while processing your command.", err
		}
	} else if screaming && len(text) > 0 {
		response = strings.ToUpper(text)
	} else {
		// Handle with Dialogflow
		handleMessageDialogflow(LINE, event, text, b)
		return "", nil // Response handled by Dialogflow
	}

	return response, nil
}

// func (b *lineBot) handleLineMessage(event *linebot.Event, message *linebot.TextMessage) {

// 	// Get user profile information from Line
// 	userProfile, err := b.lineClient.GetProfile(event.Source.UserID).Do()
// 	if err != nil {
// 		fmt.Printf("Error fetching user profile: %v\n", err)
// 		return
// 	}

// 	userID := event.Source.UserID
// 	text := message.Text

// 	fmt.Printf("User ID: %s \n", userID)

// 	//userIDInt, _ := strconv.ParseInt(userID, 10, 64)

// 	// Log the received message for debugging
// 	fmt.Printf("Received message: %s \n", text)

// 	// Check if the user exists in the database
// 	var dbUser models.User
// 	err = database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).First(&dbUser).Error

// 	// If the user does not exist, create a new user record
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			dbUser = models.User{
// 				UserID:       userID, // LINE UserID as unique identifier
// 				UserName:     userProfile.DisplayName,
// 				FirstName:    "", // LINE doesn't provide firstname and lastname
// 				LastName:     "",
// 				LanguageCode: userProfile.Language,
// 			}
// 			err = database.DB.Create(&dbUser).Error
// 			if err != nil {
// 				fmt.Printf("Error creating user: %s", err.Error())
// 				return
// 			}

// 			// Generate a JWT token for the new user
// 			token, err := token.GenerateToken(userID, "user") // Convert userID to int if needed
// 			if err != nil {
// 				fmt.Printf("Error generating JWT: %s", err.Error())
// 				return
// 			}

// 			// Send the token to the user
// 			msg := linebot.NewTextMessage("Welcome! Your access token is: " + token)
// 			if _, err := b.lineClient.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
// 				fmt.Printf("Error sending token message: %s \n", err.Error())
// 			}
// 		} else {
// 			// Handle other types of errors
// 			fmt.Printf("Error retrieving user: %s", err.Error())
// 		}
// 	} else {
// 		var response string
// 		if strings.HasPrefix(text, "/") { // Check if the message is a command by prefix "/"
// 			//response, err = handleLineCommand(event, text)
// 			response, err = handleCommand(LINE, event, text)
// 			if err != nil {
// 				fmt.Printf("An error occurred: %s \n", err.Error())
// 				response = "An error occurred while processing your command."
// 			}

// 		} else if screaming && len(text) > 0 {
// 			// If screaming mode is on, send the text in uppercase
// 			response = strings.ToUpper(text)
// 		} else {
// 			// If not a command, process the message using processLineMessage function
// 			/*response := processMessage(text)
// 			if _, err := utils.LineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
// 				fmt.Printf("An error occurred: %s \n", err.Error())
// 			}*/

// 			// Process the message using Dialogflow
// 			handleMessageDialogflow(LINE, event, text)
// 			//handleLineMessageDF(event, message)
// 			return
// 		}

// 		fmt.Printf("Response: '%s'\n", response)

// 		// Send the response if it's not empty
// 		if response != "" {
// 			if _, err := b.lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do(); err != nil {
// 				fmt.Printf("An error occurred: %s \n", err.Error())
// 			}
// 		}

// 	}

// }

// Send a message via LINE
func (b *lineBot) sendLineMessage(replyToken string, messageText string) error {
	// Create the message payload
	replyMessage := linebot.NewTextMessage(messageText)

	// Send the message
	_, err := b.lineClient.ReplyMessage(replyToken, replyMessage).Do()
	if err != nil {
		return fmt.Errorf("error sending LINE message: %w", err)
	}
	return nil
}

func (b *lineBot) sendLineMenu(replyToken string) error {
	// Define the LINE menu
	firstMenu := linebot.NewTextMessage("Here's the LINE menu:")
	actions := linebot.NewURIAction("Visit website", "https://example.com")
	template := linebot.NewButtonsTemplate("", "Menu", "Select an option:", actions)
	message := linebot.NewTemplateMessage("Menu", template)

	// Send the menu message to the user
	_, err := b.lineClient.ReplyMessage(replyToken, firstMenu, message).Do()
	if err != nil {
		return fmt.Errorf("error sending LINE menu: %w", err)
	}
	return nil
}

func (b *lineBot) ParseRequest(req *http.Request) ([]*linebot.Event, error) {
	return b.lineClient.ParseRequest(req)
}
