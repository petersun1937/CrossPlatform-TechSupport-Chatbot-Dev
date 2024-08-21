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

func DetectIntentText(projectID, sessionID, text, languageCode string) (*dialogflowpb.DetectIntentResponse, error) {
	ctx := context.Background()
	client, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
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

	return client.DetectIntent(ctx, req)
}
