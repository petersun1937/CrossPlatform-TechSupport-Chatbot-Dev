package bot

import (
	"bytes"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type TgBot interface {
	Run() error
	//SetWebhook(webhookURL string) error
	//sendTelegramMessage(chatID int64, messageText string) error
	//HandleTelegramUpdate(update tgbotapi.Update)
	//HandleTgMessage(message *tgbotapi.Message)
	//StoreDocumentChunks(filename, docID, text string, chunkSize, minchunkSize int) error
	//V2ProcessDocument(filename, fileID, filePath string) ([]Document, []string, error)
	//V2StoreDocumentChunks(filename, docID, docText string, chunkSize, overlap int) ([]Document, error)
	//ProcessDocument(filename, sessionID, filePath string) ([]Document, []string, error)
	//StoreDocumentChunks(filename, docID, chunkText string, chunkid int) (Document, error)
	GetDocFile(update tgbotapi.Update) (string, string, string, error)
	ValidateUser(user *tgbotapi.User, message *tgbotapi.Message) (bool, error)
	//sendResponse(identifier interface{}, response string) error
}

type tgBot struct {
	BaseBot
	// conf         config.BotConfig
	//embConfig config.EmbeddingConfig
	// ctx          context.Context
	// token        string
	botApi *tgbotapi.BotAPI
	//openAIclient *openai.Client
	//service *service.Service
}

// creates a new TGBot instance
func NewTGBot(botconf *config.BotConfig, embconf config.EmbeddingConfig, database database.Database, dao repository.DAO) (*tgBot, error) {
	// Attempt to create a new Telegram bot using the provided token
	botApi, err := tgbotapi.NewBotAPI(botconf.TelegramBotToken)
	if err != nil {
		return nil, err
	}
	// Ensure botApi is not nil before proceeding
	if botApi == nil {
		return nil, errors.New("telegram Bot API is nil")
	}

	return &tgBot{
		BaseBot: BaseBot{
			platform:     TELEGRAM,
			conf:         botconf,
			database:     database,
			dao:          dao,
			openAIclient: openai.NewClient(),
			embConfig:    embconf,
		},
		botApi: botApi,
		//openAIclient: openai.NewClient(),
	}, nil
}

// SetWebhook sets the webhook for Telegram bot
func (b *tgBot) setWebhook(webhookURL string) error {
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
	botApi, err := tgbotapi.NewBotAPI(b.conf.TelegramBotToken) // create new BotAPI instance using the token
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
	return b.setWebhook(b.BaseBot.conf.TelegramWebhookURL)
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
// func (b *tgBot) HandleTelegramUpdate(update tgbotapi.Update) {
// 	// Check if the update contains a message
// 	if update.Message != nil {
// 		if update.Message.Document != nil {
// 			// If the message contains a document, handle the document upload
// 			b.HandleDocumentUpload(update)
// 		} else {
// 			// Otherwise, handle regular text messages
// 			b.HandleTgMessage(update.Message)
// 		}
// 	} /*else if update.CallbackQuery != nil {
// 		// Handle button interactions (callback queries)
// 		b.handleButton(update.CallbackQuery)
// 	}*/
// }

// Handle Telegram messages
// func (b *tgBot) HandleTgMessage(message *tgbotapi.Message) {
// 	user := message.From
// 	if user == nil {
// 		return
// 	}

// 	// token, err := b.validateAndGenerateToken(userIDStr, user, message)
// 	userExists, err := b.ValidateUser(user, message)
// 	if err != nil {
// 		fmt.Printf("Error validating user: %s", err.Error())
// 		return
// 	}

// 	if !userExists {
// 		// User was just created and welcomed, no further action needed.
// 		return
// 	}

// 	// if token != nil {
// 	// 	b.sendTelegramMessage(message.Chat.ID, "Welcome! Your access token is: "+*token)
// 	// } else {
// 	b.processUserMessage(message, user.FirstName, message.Text)
// 	//}
// }

// validateUser checks if the user exists in the database and creates a new record if not.
func (b *tgBot) ValidateUser(user *tgbotapi.User, message *tgbotapi.Message) (bool, error) {
	var dbUser models.User

	userIDStr := strconv.FormatInt(user.ID, 10)
	fmt.Printf("User ID: %s \n", userIDStr)

	// Check if the user exists in the database.
	err := b.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist; create a new record.
			dbUser = models.User{
				UserID:       userIDStr,
				UserName:     user.UserName,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				LanguageCode: user.LanguageCode,
			}

			// Save the new user to the database.
			if err := b.database.GetDB().Create(&dbUser).Error; err != nil {
				return false, fmt.Errorf("error creating user: %w", err)
			}

			// Send a welcome message to the new user.
			welcomeMessage := fmt.Sprintf("Welcome, %s!", user.UserName)
			if err := b.SendReply(message, welcomeMessage); err != nil {
				return false, fmt.Errorf("error sending welcome message: %w", err)
			}

			return false, nil // User was just created and welcomed.
		}
		return false, fmt.Errorf("error retrieving user: %w", err)
	}
	return true, nil // User already exists.
}

// validateAndGenerateToken checks if the user exists in the database and generates a token if not
// func (b *tgBot) validateAndGenerateToken(userIDStr string, user *tgbotapi.User) (*string, error) {
// 	// Check if the user exists in the database
// 	var dbUser models.User
// 	err := b.Service.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
// 	if err != nil {
// 		// If user does not exist, create a new user
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			dbUser = models.User{
// 				UserID:       userIDStr,
// 				UserName:     user.UserName,
// 				FirstName:    user.FirstName,
// 				LastName:     user.LastName,
// 				LanguageCode: user.LanguageCode,
// 			}

// 			// Create the new user record in the database
// 			if err := b.Service.GetDB().Create(&dbUser).Error; err != nil {
// 				return nil, fmt.Errorf("error creating user: %w", err)
// 			}

// 			// Generate a JWT token using the service's ValidateUser method
// 			token, err := b.Service.ValidateUser(userIDStr, service.ValidateUserReq{
// 				FirstName:    user.FirstName,
// 				LastName:     user.LastName,
// 				UserName:     user.UserName,
// 				LanguageCode: user.LanguageCode,
// 			})
// 			if err != nil {
// 				return nil, fmt.Errorf("error generating JWT: %w", err)
// 			}

// 			return token, nil // Return the generated token
// 		}
// 		return nil, fmt.Errorf("error retrieving user: %w", err)
// 	}

// 	return nil, nil // User already exists, no token generation needed
// }

// Process user messages and respond accordingly
// func (b *tgBot) processUserMessage(message *tgbotapi.Message, firstName, text string) { //chatID int64
// 	chatID := message.Chat.ID

// 	fmt.Printf("Received message from %s: %s \n", firstName, text)
// 	fmt.Printf("Chat ID: %d \n", chatID)

// 	var response string
// 	//var err error

// 	if strings.HasPrefix(text, "/") {
// 		response = b.BaseBot.HandleCommand(text)
// 		/*response, err = handleCommand(chatID, text, b)
// 		if err != nil {
// 			fmt.Printf("An error occurred: %s \n", err.Error())
// 			response = "An error occurred while processing your command."
// 		}*/
// 	} else if b.conf.Screaming && len(text) > 0 {
// 		response = strings.ToUpper(text)
// 	} else {
// 		// Get all document embeddings
// 		documentEmbeddings, chunkText, err := b.BaseBot.dao.FetchEmbeddings()
// 		//documentEmbeddings, chunkText, err := b.Service.GetAllDocumentEmbeddings()
// 		if err != nil {
// 			fmt.Printf("Error retrieving document embeddings: %v", err)
// 			response = "Error retrieving document embeddings."
// 		} else if b.conf.UseOpenAI {
// 			//conf := config.GetConfig()
// 			// Perform similarity matching with the user's message when OpenAI is enabled
// 			topChunksText, err := document.RetrieveTopNChunks(text, documentEmbeddings, b.embConfig.NumTopChunks, chunkText, b.embConfig.ScoreThreshold) // Returns maximum N chunks with similarity threshold
// 			if err != nil {
// 				fmt.Printf("Error retrieving document chunks: %v", err)
// 				response = "Error retrieving related document information."
// 			} else if len(topChunksText) > 0 {
// 				// If there are similar chunks found, respond with those
// 				//response = fmt.Sprintf("Found related information:\n%s", strings.Join(topChunksText, "\n"))

// 				// If there are similar chunks found, provide them as context for GPT
// 				context := strings.Join(topChunksText, "\n")
// 				gptPrompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, text)

// 				// Call GPT with the context and user query
// 				response, err = b.BaseBot.GetOpenAIResponse(gptPrompt)
// 				if err != nil {
// 					response = fmt.Sprintf("OpenAI Error: %v", err)
// 				} /*else {
// 					response = fmt.Sprintf("Found related information based on context:\n%s", response)
// 				}*/
// 			} else {
// 				// If no relevant document found, fallback to OpenAI response
// 				response, err = b.BaseBot.GetOpenAIResponse(text)
// 				//response = fmt.Sprintf("Found related information:\n%s", strings.Join(topChunksText, "\n"))
// 				if err != nil {
// 					response = fmt.Sprintf("OpenAI Error: %v", err)
// 				}
// 			}
// 		} else {
// 			// Fall back to Dialogflow if OpenAI is not enabled
// 			//b.BaseBot.handleMessageDialogflow(TELEGRAM, message, text, b) //TODO
// 			return
// 		}
// 	}

// 	/*if strings.HasPrefix(text, "/") {
// 		response, err = handleCommand(chatID, text, b)
// 		if err != nil {
// 			fmt.Printf("An error occurred: %s \n", err.Error())
// 			response = "An error occurred while processing your command."
// 		}
// 	} else if screaming && len(text) > 0 {
// 		response = strings.ToUpper(text)
// 	} else if useOpenAI {
// 		// Call OpenAI to get the response
// 		response, err = GetOpenAIResponse(text)
// 		if err != nil {
// 			response = fmt.Sprintf("OpenAI Error: %v", err)
// 		}
// 	} else {
// 		handleMessageDialogflow(TELEGRAM, message, text, b)
// 		return
// 	}*/

// 	if response != "" {
// 		b.sendTelegramMessage(chatID, response)
// 	}
// }

// handleDialogflowResponse processes and sends the Dialogflow response to the appropriate platform
// func (b *tgBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {

// 	// Send the response to respective platform
// 	// by iterating over the fulfillment messages returned by Dialogflow and processes any text messages.
// 	for _, msg := range response.QueryResult.FulfillmentMessages {
// 		if _, ok := identifier.(*tgbotapi.Message); ok {
// 			if text := msg.GetText(); text != nil {
// 				return b.SendReply(identifier, text.Text[0])
// 			}
// 		}
// 	}
// 	return fmt.Errorf("invalid Telegram message identifier")
// }

// Check identifier and send message via Telegram
func (b *tgBot) SendReply(identifier interface{}, response string) error {
	if message, ok := identifier.(*tgbotapi.Message); ok { // Assertion to check if identifier is of type tgbotapi.Message
		return b.sendTelegramMessage(message.Chat.ID, response)
	} else {
		return fmt.Errorf("invalid identifier for Telegram platform")
	}
}

// Send a message via Telegram (TG requires manual construction of an HTTP request)
func (b *tgBot) sendTelegramMessage(chatID int64, messageText string) error {
	// Use the Telegram API URL from the config
	url := b.conf.TelegramAPIURL + b.conf.TelegramBotToken + "/sendMessage"
	//conf := config.GetConfig()
	//url := conf.TelegramAPIURL + b.token + "/sendMessage"

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

// GetDocFile retrieves relevant file information of the document
func (b *tgBot) GetDocFile(update tgbotapi.Update) (string, string, string, error) {

	fileID := update.Message.Document.FileID
	filename := update.Message.Document.FileName
	fileURL, err := b.botApi.GetFileDirectURL(fileID)
	if err != nil {
		//b.sendTelegramMessage(update.Message.Chat.ID, "Error getting file: "+err.Error())
		return "", "", "", err
	}
	return fileID, fileURL, filename, nil
}

/*func (b *tgBot) ProcessDocument(filename, sessionID, filePath string) ([]Document, []string, error) {
	// Extract text from the uploaded file (assuming downloadAndExtractText can handle local files)
	docText, err := document.DownloadAndExtractText(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing document: %w", err)
	}

	chunks := document.OverlapChunk(docText, b.embConfig.ChunkSize, b.embConfig.MinChunkSize)

	documents := make([]Document, 0)

	tagList := []string{} // Initialize the tag list

	for i, chunk := range chunks {

		document, err := b.storeDocumentChunks(filename, filename+"_"+sessionID, chunk, i)
		if err != nil {
			return nil, nil, err
		}

		documents = append(documents, document)

		// Auto-tagging using OpenAI

		tags, err := b.openAIclient.AutoTagWithOpenAI(docText)
		if err != nil {
			return nil, nil, fmt.Errorf("error auto-tagging document: %w", err)
		}

		// Append tags to the tag list
		tagList = append(tagList, tags...)
	}

	// Remove duplicates from the tag list
	uniqueTags := utils.RemoveDuplicates(tagList)

	return documents, uniqueTags, nil
}

func (b *tgBot) storeDocumentChunks(filename, docID, chunkText string, chunkid int) (Document, error) {
	// Chunk the document with overlap

	//client := openai.NewClient()

	//for i, chunk := range chunks {
	// Get the embeddings for each chunk
	embedding, err := b.openAIclient.EmbedText(chunkText)
	if err != nil {
		return Document{}, fmt.Errorf("error embedding chunk %d: %v", chunkid, err)
	}

	// Create a unique chunk ID for storage in the database
	chunkID := fmt.Sprintf("%s_chunk_%d_%s", filename, chunkid, docID)

	chunkText = utils.SanitizeText(chunkText)
	embeddingStr := utils.Float64SliceToPostgresArray(embedding)

	document := Document{
		Filename: filename,
		DocID:    docID,
		ChunkID:  chunkID,
		//DocText:   docText,
		DocText:   chunkText,
		Embedding: embeddingStr,
	}

	//}

	return document, nil
}*/

/*
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
*/
/*
func (b *tgBot) V2ProcessDocument(filename, fileID, filePath string) ([]Document, []string, error) {
	// Extract text from the uploaded file (assuming downloadAndExtractText can handle local files)
	docText, err := document.DownloadAndExtractText(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing document: %w", err)
	}

	documents, err := b.V2StoreDocumentChunks(filename, filename+"_"+fileID, docText, b.embConfig.ChunkSize, b.embConfig.MinChunkSize)
	if err != nil {
		return nil, nil, err
	}

	// Auto-tagging using OpenAI
	tags, err := b.openAIclient.AutoTagWithOpenAI(docText)
	if err != nil {
		return nil, nil, fmt.Errorf("error auto-tagging document: %w", err)
	}

	return documents, tags, nil
}

func (b *tgBot) V2StoreDocumentChunks(filename, docID, docText string, chunkSize, overlap int) ([]Document, error) {
	// Chunk the document with overlap
	chunks := document.OverlapChunk(docText, chunkSize, overlap)

	//client := openai.NewClient()

	documents := make([]Document, 0)
	for i, chunk := range chunks {
		// Get the embeddings for each chunk
		embedding, err := b.openAIclient.EmbedText(chunk)
		if err != nil {
			return nil, fmt.Errorf("error embedding chunk %d: %v", i, err)
		}

		// Create a unique chunk ID for storage in the database
		chunkID := fmt.Sprintf("%s_chunk_%d_%s", filename, i, docID)

		docText = utils.SanitizeText(docText)
		embeddingStr := utils.Float64SliceToPostgresArray(embedding)

		document := Document{
			Filename:  filename,
			DocID:     docID,
			ChunkID:   chunkID,
			DocText:   chunk,
			Embedding: embeddingStr,
		}
		documents = append(documents, document)
	}

	return documents, nil
}
*/
