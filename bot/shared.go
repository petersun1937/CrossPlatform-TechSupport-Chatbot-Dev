package bot

import (
	"Tg_chatbot/utils"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// For screaming state
var screaming bool

// Process the commands sent by users and returns the message as a string
func handleCommand(identifier interface{}, command string, bot Bot) (string, error) {
	var message string
	var err error

	switch command {
	case "/start":
		message = "Welcome to the bot!"
	case "/scream":
		screaming = true // Enable screaming mode
		message = "Scream mode enabled!"
	case "/whisper":
		screaming = false // Disable screaming mode
		message = "Scream mode disabled!"
	case "/menu":
		// Handle menu sending based on platform
		/*switch platform {
		case LINE:
			if event, ok := identifier.(*linebot.Event); ok {
				err = sendLineMenu(event.ReplyToken) // Send a menu to LINE
			} else {
				err = fmt.Errorf("invalid identifier type for LINE platform")
			}
		case TELEGRAM:
			if chatID, ok := identifier.(int64); ok {
				err = sendMenu(chatID) // Send a menu to Telegram
			} else {
				err = fmt.Errorf("invalid identifier type for Telegram platform")
			}
		}*/
		err = bot.sendMenu(identifier)
		if err != nil {
			return "", err
		}
		return "", nil
	case "/help":
		message = "Here are some commands you can use: /start, /help, /scream, /whisper, /menu"
	case "/custom":
		message = "This is a custom response!"
	default:
		message = "I don't know that command"
	}

	return message, nil
}

// handleMessageDialogflow handles messages from both LINE and Telegram
func handleMessageDialogflow(platform Platform, identifier interface{}, text string, bot Bot) {
	// Determine sessionID based on platform
	sessionID, err := getSessionID(platform, identifier)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send the message to Dialogflow and receive a response TODO: testagent-mkyg
	response, err := utils.DetectIntentText("testagent-mkyg", sessionID, text, "en")
	if err != nil {
		fmt.Printf("Error detecting intent: %v\n", err)
		return
	}

	// Process and send the Dialogflow response to the appropriate platform
	if err := bot.handleDialogflowResponse(response, identifier); err != nil {
		fmt.Println(err)
	}
}

// getSessionID extracts the session ID based on the platform and identifier
func getSessionID(platform Platform, identifier interface{}) (string, error) {
	switch platform {
	case LINE:
		if event, ok := identifier.(*linebot.Event); ok {
			return event.Source.UserID, nil
		}
		return "", fmt.Errorf("invalid LINE event identifier")
	case TELEGRAM:
		if message, ok := identifier.(*tgbotapi.Message); ok {
			return strconv.FormatInt(message.Chat.ID, 10), nil
		}
		return "", fmt.Errorf("invalid Telegram message identifier")
	case GENERAL:
		if sessionID, ok := identifier.(string); ok {
			return sessionID, nil
		}
		return "", fmt.Errorf("invalid Telegram message identifier")
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}

/*
type Platform int

const (
	LINE Platform = iota
	TELEGRAM
)

// For screaming state
var screaming bool*/

var keywords = []string{
	"hello",
	"help",
	"start",
	"menu",
	"scream",
	"whisper",
}

// Define responses for different messages
func processMessage(message string) string {
	// Convert message to lowercase to ensure case-insensitive matching
	message = strings.ToLower(message)

	// Perform fuzzy matching
	bestMatch := fuzzy.RankFind(message, keywords)

	if len(bestMatch) > 0 {
		switch bestMatch[0].Target {
		case "hello":
			return "Hello! How can I assist you today?"
		case "help":
			return "Here are some commands you can use: /start, /help, /scream, /whisper, /menu"
		case "start":
			return "Let's get started!"
		case "menu":
			return "Here's the menu: ..."
		case "scream":
			screaming = true
			return "Scream mode enabled! (Type /whisper to disable)"
		/*case "whisper":
		utils.Screaming = false
		return "Whisper mode enabled!"*/
		default:
			return "I'm sorry, I didn't understand that. Type /help to see what I can do."
		}
	}

	return "I'm sorry, I didn't understand that. Type /help to see what I can do."
}

/*func (b *BaseBot) sendResponse(identifier interface{}, response string) error { // identifier is chatID for TG, reply token for LINE
	switch b.Platform {
	case LINE: // If the platform is LINE
		if event, ok := identifier.(*linebot.Event); ok { // assertion to check if identifier is of type linebot.Event
			return b.sendLineMessage(event.ReplyToken, response)

		} else {
			return fmt.Errorf("Invalid identifier for LINE platform")
		}
	case TELEGRAM: // If the platform is Telegram
		if chatID, ok := identifier.(int64); ok { // assertion to check if identifier is of type int64
			return b.sendTelegramMessage(chatID, response)
		} else {
			return fmt.Errorf("Invalid identifier for Telegram platform")
		}
	default:
		return fmt.Errorf("Unsupported platform")
	}
}*/
