package handlers

import (
	"Tg_chatbot/utils"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Platform int

const (
	LINE Platform = iota
	TELEGRAM
)

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

// Handle messages with Dialogflow
func handleMessageDialogflow(message *tgbotapi.Message) {
	projectID := "testagent-mkyg"
	sessionID := strconv.FormatInt(message.Chat.ID, 10)
	languageCode := "en" // or your preferred language

	response, err := utils.DetectIntentText(projectID, sessionID, message.Text, languageCode)
	if err != nil {
		fmt.Printf("Error detecting intent: %v\n", err)
		return
	}

	handleDialogflowResponse(response, TELEGRAM, message.Chat.ID)
}

func handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, platform Platform, identifier interface{}) {
	if response != nil {
		for _, msg := range response.QueryResult.FulfillmentMessages {
			if text := msg.GetText(); text != nil {
				switch platform {
				case LINE:
					if event, ok := identifier.(*linebot.Event); ok {
						replyMessage := linebot.NewTextMessage(text.Text[0])
						_, err := utils.LineBot.ReplyMessage(event.ReplyToken, replyMessage).Do()
						if err != nil {
							fmt.Printf("Error sending LINE reply: %v\n", err)
						}
					}
				case TELEGRAM:
					if chatID, ok := identifier.(int64); ok {
						utils.SendTelegramResponse(chatID, text.Text[0])
					}
				}
			}
		}
	}
}
