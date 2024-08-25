package bot

import (
	"Tg_chatbot/config"
	"Tg_chatbot/service"
	"Tg_chatbot/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgBot struct {
	ctx     context.Context
	token   string
	botApi  *tgbotapi.BotAPI
	service *service.Service
}

func NewTGBot(conf *config.Config, service *service.Service) Bot {
	return &tgBot{
		ctx:     context.Background(),
		token:   conf.GetTelegramBotToken(),
		service: service,
	}
}

func (b *tgBot) Run() error {
	botApi, err := tgbotapi.NewBotAPI(b.token) // create new BotAPI instance using the token
	// utils.TgBot: global variable (defined in the utils package) that holds the reference to the bot instance.
	if err != nil {
		return err
	}

	b.botApi = botApi

	// FIXME: ???
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // timeout for long polling set to 60 s

	// Get updates channel to start long polling to receive updates.
	// The channel will be continuously fed with new Update objects from Telegram.
	updates := b.botApi.GetUpdatesChan(u)

	// Use go routine to continuously process received updates from the updates channel
	go b.receiveUpdates(b.ctx, updates)
	return nil
}

// ReceiveUpdates receives updates from Telegram API and handles them
func (b *tgBot) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// "updates" is a channel that receives updates from the Telegram bot (e.g., messages, button clicks).
	// The bot's API sends these updates to the application, and the function processes them by handling the updates.

	for { // continuous loop to check for updates
		select { // select statement waits for one of its cases to be ready, then executes the first case that becomes available.
		case <-ctx.Done(): // if context has been cancelled:
			fmt.Println("Goroutine: Received cancel signal, stopping...")
			// exit the loop and stop the go routine
			return
		case update := <-updates: // Process incoming updates from Telegram
			b.HandleTelegramUpdate(update)
		}
	}
}

// HandleTelegramUpdate processes incoming updates from Telegram
func (b *tgBot) HandleTelegramUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		b.handleTgMessage(update.Message) // handle user input message
	} else if update.CallbackQuery != nil { // a callback query is typically generated when a user interacts with an inline button within a message.
		handleButton(update.CallbackQuery) // handle button press activated by sendMenu
	}
}

// Processes incoming messages from users
func (b *tgBot) handleTgMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Convert user.ID from int64 to string
	userIDStr := strconv.FormatInt(user.ID, 10)

	fmt.Printf("User ID: %s \n", userIDStr)

	req := service.ValidateUserReq{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserName:     user.UserName,
		LanguageCode: user.LanguageCode,
	}
	token, err := b.service.ValidateUser(userIDStr, req)
	if err != nil {
		// TODO: log error
		return
	}

	// Send the token to the user
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Your access token is: "+*token)
	utils.TgBot.Send(msg)

	// else {
	// 	fmt.Printf("Received message from %s: %s \n", user.FirstName, text)
	// 	fmt.Printf("Chat ID: %d \n", message.Chat.ID)

	// 	var response string
	// 	if strings.HasPrefix(text, "/") {
	// 		// Handle commands
	// 		//response, err = handleCommand(message.Chat.ID, text)
	// 		response, err = handleCommand(TELEGRAM, message.Chat.ID, text)
	// 		if err != nil {
	// 			fmt.Printf("An error occurred: %s \n", err.Error())
	// 			response = "An error occurred while processing your command."
	// 		}
	// 	} else if screaming && len(text) > 0 {
	// 		// If screaming mode is on, send the text in uppercase
	// 		response = strings.ToUpper(text)
	// 	} else {
	// 		// Process the message using processMessage function
	// 		//response = processMessage(text)

	// 		// Send the message to Dialogflow for processing
	// 		handleMessageDialogflow(TELEGRAM, message.Chat.ID, text)
	// 		//handleTGMessageDialogflow(message)
	// 		return
	// 	}

	// 	fmt.Printf("Response: '%s'\n", response)

	// 	// Send the response if it's not empty
	// 	if response != "" {
	// 		msg := tgbotapi.NewMessage(message.Chat.ID, response)
	// 		_, err = utils.TgBot.Send(msg)
	// 		if err != nil {
	// 			fmt.Printf("An error occurred: %s \n", err.Error())
	// 		}
	// 	}
	// }
}

// Handle messages with Dialogflow
/*func handleTGMessageDialogflow(message *tgbotapi.Message) {
	projectID := "testagent-mkyg"
	sessionID := strconv.FormatInt(message.Chat.ID, 10)
	languageCode := "en"

	// Send the user’s message to Dialogflow and receives a response.
	response, err := utils.DetectIntentText(projectID, sessionID, message.Text, languageCode)
	if err != nil {
		fmt.Printf("Error detecting intent: %v\n", err)
		return
	}

	// Process Dialogflow response and send it
	handleDialogflowResponse(response, TELEGRAM, message.Chat.ID)
}*/

func handleButton(query *tgbotapi.CallbackQuery) {
	var text string
	markup := tgbotapi.NewInlineKeyboardMarkup()

	if query.Data == utils.NextButton {
		text = utils.SecondMenu
		markup = utils.SecondMenuMarkup
	} else if query.Data == utils.BackButton {
		text = utils.FirstMenu
		markup = utils.FirstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	utils.TgBot.Send(callbackCfg)

	msg := tgbotapi.NewEditMessageTextAndMarkup(query.Message.Chat.ID, query.Message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	utils.TgBot.Send(msg)
}

// processes custom messages

func HandleCustomMessage(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}
	// Bind the incoming JSON payload to the request struct
	if err := c.BindJSON(&req); err != nil {
		fmt.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	fmt.Printf("Received message from frontend: %s", req.Message)
	// Process the message and generate a response
	response := processMessage(req.Message)
	c.JSON(http.StatusOK, gin.H{"response": response})
}

/*
// generates a response for the given message
func processMessage(message string) string {
	// Placeholder for NLP model integration
	return "Processed message: " + message
}*/
