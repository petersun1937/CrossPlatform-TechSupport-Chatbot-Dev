package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Global variables to hold the bot instance and screaming state
var TgBot *tgbotapi.BotAPI
var Screaming bool

// Telegram API URL
const telegramAPIURL = "https://api.telegram.org/bot"

// SendTelegramResponse sends a response to a Telegram chat
/*
func SendTelegramResponse(chatID int64, response string) {
	// Construct the URL for the Telegram API request
	url := telegramAPIURL + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	// Create the message payload
	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    response,
	}

	// Marshal the message payload to JSON
	jsonMessage, _ := json.Marshal(message)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}
	defer resp.Body.Close()

	// Log the response
	log.Printf("Response sent to chat ID %d", chatID)
}*/

// Send a message via LINE
func SendLineMessage(replyToken string, messageText string) error {
	replyMessage := linebot.NewTextMessage(messageText)
	_, err := LineBot.ReplyMessage(replyToken, replyMessage).Do()
	return err
}

// Send a message via Telegram
func SendTelegramMessage(chatID int64, messageText string) error {
	// Construct the URL for the Telegram API request
	url := telegramAPIURL + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	// Create the message payload
	message := map[string]interface{}{
		"chat_id": chatID,
		"text":    messageText,
	}

	// Marshal the message payload to JSON
	jsonMessage, _ := json.Marshal(message)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	defer resp.Body.Close()

	// Log the response (can be removed if not needed)
	log.Printf("Response sent to chat ID %d", chatID)

	return nil
}

// Menu texts
var (
	FirstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	SecondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	NextButton     = "Next"
	BackButton     = "Back"
	TutorialButton = "Tutorial"

	// Keyboard layout for the first menu. One button, one row
	FirstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(NextButton, NextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	SecondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BackButton, BackButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(TutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

// SendMenu sends a menu to the chat
func SendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, FirstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = FirstMenuMarkup
	_, err := TgBot.Send(msg)
	return err
}

// For Line bot
var LineBot *linebot.Client

func InitLineBot(channelSecret, channelToken string) error {
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return err
	}
	LineBot = bot
	return nil
}

func SendLineMenu(replyToken string) error {
	// ***
	return nil
}

// Send a text query to Dialogflow and returns the response
func DetectIntentText(projectID, sessionID, text, languageCode string) (*dialogflowpb.DetectIntentResponse, error) {
	// Create a background context for the API call
	ctx := context.Background()

	// Create a new Dialogflow Sessions client
	client, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close() // Ensure the client is closed when done

	// Construct the session path for the Dialogflow API
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	// Create the DetectIntentRequest with the session path and query input
	req := &dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         text,
					LanguageCode: languageCode,
				},
			},
		},
	}

	// Send the request and return the response or error
	return client.DetectIntent(ctx, req)
}
