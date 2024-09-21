package bot

import (
	config "Tg_chatbot/configs"
	"Tg_chatbot/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type TgBot interface {
	Run() error
	SetWebhook(webhookURL string) error
	HandleTelegramUpdate(update tgbotapi.Update)
}

type tgBot struct {
	BaseBot
	conf   config.BotConfig
	ctx    context.Context
	token  string
	botApi *tgbotapi.BotAPI
	//service *service.Service
}

func NewTGBot(conf *config.Config, service service.Service) (*tgBot, error) {
	// Attempt to create a new Telegram bot using the provided token
	botApi, err := tgbotapi.NewBotAPI(conf.TelegramBotToken)
	if err != nil {
		return nil, err
	}

	// Ensure botApi is not nil before proceeding
	if botApi == nil {
		return nil, errors.New("telegram Bot API is nil")
	}

	baseBot := BaseBot{
		Platform: TELEGRAM,
		Service:  service,
	}

	// Initialize and return tgBot instance
	return &tgBot{
		conf:    conf.BotConfig,
		BaseBot: baseBot,
		ctx:     context.Background(),
		token:   conf.TelegramBotToken,
		botApi:  botApi,
	}, nil

}

// SetWebhook sets the webhook for Telegram bot
func (b *tgBot) sSetWebhook(webhookURL string) error {
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return fmt.Errorf("error creating webhook config: %w", err)
	}

	// Send the request to set the webhook
	_, err = b.botApi.Request(webhookConfig)
	if err != nil {
		return fmt.Errorf("error setting webhook: %w", err)
	}

	return nil
}

func (b *tgBot) Run() error {
	botApi, err := tgbotapi.NewBotAPI(b.token) // create new BotAPI instance using the token
	// utils.TgBot: global variable (defined in the utils package) that holds the reference to the bot instance.
	if err != nil {
		return err
	}

	b.botApi = botApi

	// Start the bot with webhook
	fmt.Println("Telegram bot is running with webhook!")

	// /// Create a new update configuration with offset of 0
	// // Using 0 means it will start fetching updates from the beginning.
	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60 // timeout for long polling set to 60 s

	// // Get updates channel to start long polling to receive updates.
	// // The channel will be continuously fed with new Update objects from Telegram.
	// updates := b.botApi.GetUpdatesChan(u)

	// // Use go routine to continuously process received updates from the updates channel
	// go b.receiveUpdates(b.ctx, updates)
	return b.sSetWebhook(b.conf.GetTelegramBotToken())
}

// Receives updates from Telegram API and handles them (for long polling, not needed with Webhook)
// func (b *tgBot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
// 	// "updates" is a channel that receives updates from the Telegram bot (e.g., messages, button clicks).
// 	// The bot's API sends these updates to the application, and the function processes them by handling the updates.

// 	for { // continuous loop to check for updates
// 		select { // select statement waits for one of its cases to be ready, then executes the first case that becomes available.
// 		case <-ctx.Done(): // if context has been cancelled:
// 			fmt.Println("Goroutine: Received cancel signal, stopping...")
// 			// exit the loop and stop the go routine
// 			return
// 		case update := <-updates: // Process incoming updates from Telegram
// 			b.HandleTelegramUpdate(update)
// 		}
// 	}
// }

// HandleTelegramUpdate processes incoming updates from Telegram
func (b *tgBot) HandleTelegramUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		b.handleTgMessage(update.Message) // handle user input message
	} else if update.CallbackQuery != nil { // a callback query is typically generated when a user interacts with an inline button within a message.
		b.handleButton(update.CallbackQuery) // handle button press activated by sendMenu
	}
}

func (b *tgBot) GetUserID() (*string, error) {
	userProfile, err := b.lineClient.GetProfile(userID).Do()
	if err != nil {
		return nil, err
	}
	return &userProfile.UserID, nil
}

// Handle Telegram messages
func (b *tgBot) handleTgMessage(message *tgbotapi.Message) {
	user := message.From
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
	}
}

// Validate user in the database and generate token if user is new
func (b *tgBot) validateAndGenerateToken(userIDStr string, user *tgbotapi.User) (*string, error) {
	req := service.ValidateUserReq{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		LanguageCode: user.LanguageCode,
	}

	token, err := b.Service.ValidateUser(userIDStr, req)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("User not found, created record")
		return token, nil
	}
	return token, err
}

// Process user messages and respond accordingly
func (b *tgBot) processUserMessage(message *tgbotapi.Message, firstName, text string) { //chatID int64
	chatID := message.Chat.ID

	fmt.Printf("Received message from %s: %s \n", firstName, text)
	fmt.Printf("Chat ID: %d \n", chatID)

	var response string
	var err error

	if strings.HasPrefix(text, "/") {
		response, err = handleCommand(chatID, text, b)
		if err != nil {
			fmt.Printf("An error occurred: %s \n", err.Error())
			response = "An error occurred while processing your command."
		}
	} else if screaming && len(text) > 0 {
		response = strings.ToUpper(text)
	} else {
		handleMessageDialogflow(TELEGRAM, message, text, b)
		return
	}

	if response != "" {
		b.sendTelegramMessage(chatID, response)
	}
}

// handleDialogflowResponse processes and sends the Dialogflow response to the appropriate platform
func (b *tgBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {

	// Send the response to respective platform
	// by iterating over the fulfillment messages returned by Dialogflow and processes any text messages.
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if _, ok := identifier.(*tgbotapi.Message); ok {
			if text := msg.GetText(); text != nil {
				return b.sendResponse(identifier, text.Text[0])
			}
		}
	}
	return fmt.Errorf("invalid Telegram message identifier")
}

// Check identifier and send message via Telegram
func (b *tgBot) sendResponse(identifier interface{}, response string) error {
	if message, ok := identifier.(*tgbotapi.Message); ok { // Assertion to check if identifier is of type tgbotapi.Message
		return b.sendTelegramMessage(message.Chat.ID, response)
	} else {
		return fmt.Errorf("invalid identifier for Telegram platform")
	}
}

// Send a message via Telegram (TG requires manual construction of an HTTP request)
func (b *tgBot) sendTelegramMessage(chatID int64, messageText string) error {
	// Use the Telegram API URL from the config
	conf := config.GetConfig()
	url := conf.TelegramAPIURL + b.token + "/sendMessage"

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
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"

	// Keyboard layout for the first menu. One button, one row
	FirstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	SecondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

func (b *tgBot) handleButton(query *tgbotapi.CallbackQuery) {
	var text string
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if query.Data == nextButton {
		text = secondMenu
		markup = SecondMenuMarkup
	} else if query.Data == backButton {
		text = firstMenu
		markup = FirstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	b.botApi.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	b.botApi.Send(msg)
}

func (b *tgBot) sendMenu(identifier interface{}) error {
	if chatID, ok := identifier.(int64); ok {
		return b.sendTGMenu(chatID)
	} else {
		return fmt.Errorf("invalid identifier type for Telegram platform")
	}
}

// Send a menu to the Telegram chat
func (b *tgBot) sendTGMenu(chatID int64) error {
	// Define the Telegram menu
	firstMenu := "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	msg := tgbotapi.NewMessage(chatID, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = FirstMenuMarkup

	_, err := b.botApi.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending Telegram menu: %w", err)
	}
	return nil
}
