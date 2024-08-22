package handlers

import (
	"Tg_chatbot/utils"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Platform int

const (
	LINE Platform = iota
	TELEGRAM
)

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
			utils.Screaming = true
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

// Process the commands sent by users and returns the message as a string
func handleCommand(platform Platform, identifier interface{}, command string) (string, error) {
	var message string
	var err error

	switch command {
	case "/start":
		message = "Welcome to the bot!"
	case "/scream":
		utils.Screaming = true // Enable screaming mode
		message = "Scream mode enabled!"
	case "/whisper":
		utils.Screaming = false // Disable screaming mode
		message = "Scream mode disabled!"
	case "/menu":
		// Handle menu sending based on platform
		switch platform {
		case LINE:
			if event, ok := identifier.(*linebot.Event); ok {
				err = utils.SendLineMenu(event.ReplyToken) // Send a menu to LINE
			} else {
				err = fmt.Errorf("invalid identifier type for LINE platform")
			}
		case TELEGRAM:
			if chatID, ok := identifier.(int64); ok {
				err = utils.SendMenu(chatID) // Send a menu to Telegram
			} else {
				err = fmt.Errorf("invalid identifier type for Telegram platform")
			}
		}
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

// Handle responses from Dialogflow
func handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, platform Platform, identifier interface{}) {

	if response == nil {
		return
	}

	// Send the response to respective platform
	// by iterating over the fulfillment messages returned by Dialogflow and processes any text messages.
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if text := msg.GetText(); text != nil {
			sendResponse(platform, identifier, text.Text[0])
		}
	}

}

// SendResponse sends a message to the specified platform
func sendResponse(platform Platform, identifier interface{}, response string) { // identifier is chatID for TG, reply token for LINE
	switch platform {
	case LINE: // If the platform is LINE
		if event, ok := identifier.(*linebot.Event); ok { // assertion to check if identifier is of type linebot.Event
			if err := utils.SendLineMessage(event.ReplyToken, response); err != nil {
				log.Printf("Error sending LINE response: %v", err)
			}
		} else {
			log.Printf("Invalid identifier for LINE platform")
		}
	case TELEGRAM: // If the platform is Telegram
		if chatID, ok := identifier.(int64); ok { // assertion to check if identifier is of type int64
			if err := utils.SendTelegramMessage(chatID, response); err != nil {
				log.Printf("Error sending Telegram response: %v", err)
			}
		} else {
			log.Printf("Invalid identifier for Telegram platform")
		}
	default:
		log.Printf("Unsupported platform")
	}
}
