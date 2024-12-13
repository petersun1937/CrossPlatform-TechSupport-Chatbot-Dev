package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"
	"fmt"
	"net/http"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
)

type Document struct {
	Filename string
	DocID    string
	ChunkID  string
	DocText  string
	//Embedding []float64 `gorm:"type:float8[]"`
	Embedding string `gorm:"type:float8[]"` // Store as a string and ensure it's passed correctly
}

type GeneralBot interface {
	Run() error
	//HandleGeneralMessage(sessionID, message string)
	//SendResponse(identifier interface{}, response string) error
	//StoreDocumentChunks(Filename, docID, text string, chunkSize, minchunkSize int) error
	//ProcessDocument(Filename, sessionID, filePath string) error
	//V2ProcessDocument(Filename, sessionID, filePath string) ([]Document, []string, error)
	//V2StoreDocumentChunks(filename, docID, docText string, chunkSize, overlap int) ([]Document, error)
	//ProcessDocument(filename, sessionID, filePath string) ([]Document, []string, error)
	//StoreDocumentChunks(filename, docID, chunkText string, chunkid int) (Document, error)
	StoreContext(sessionID string, c *gin.Context)
	//SetWebhook(webhookURL string) error
}

type generalBot struct {
	// Add any common fields if necessary, like configuration
	BaseBot
	//ctx context.Context
	// conf         config.BotConfig
	//embConfig    config.EmbeddingConfig
	//openAIclient *openai.Client
	//config map[string]string
}

// func NewGeneralBot(conf *config.Config, service *service.Service) (*generalBot, error) {
// 	baseBot := &BaseBot{
// 		Platform: GENERAL,
// 		Service:  service,
// 	}

// 	return &generalBot{
// 		BaseBot:      baseBot,
// 		conf:         conf.BotConfig,
// 		embConfig:    conf.EmbeddingConfig,
// 		ctx:          context.Background(),
// 		openAIclient: openai.NewClient(),
// 	}, nil
// }

// creates a new GeneralBot instance
func NewGeneralBot(botconf *config.BotConfig, embconf config.EmbeddingConfig, database database.Database, dao repository.DAO) (*generalBot, error) {

	return &generalBot{
		BaseBot: BaseBot{
			platform:     GENERAL,
			conf:         botconf,
			database:     database,
			dao:          dao,
			openAIclient: openai.NewClient(),
			embConfig:    embconf,
		},
	}, nil
}

func (b *generalBot) Run() error {
	// Implement logic for running the bot
	fmt.Println("General bot is running...")
	return nil
}

// func (b *generalBot) HandleGeneralMessage(c *gin.Context) {
// func (b *generalBot) HandleGeneralMessage(sessionID, message string) {

// 	// Process and send the message
// 	b.ProcessUserMessage(sessionID, message)

// 	// Send the response back to the frontend using sendResponse
// 	/*err = b.sendFrontendMessage(c, response)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
// 		return
// 	}*/

// 	/*if err := b.sendResponse(req.SessionID, response); err != nil {
// 		fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
// 		return
// 	}*/
// }

// ProcessUserMessage processes incoming messages
// func (b *generalBot) ProcessUserMessage(sessionID string, message string) {
// 	var response string
// 	//var err error

// 	fmt.Printf("Received message %s \n", message)
// 	fmt.Printf("Chat ID: %s \n", sessionID)

// 	if strings.HasPrefix(message, "/") {

// 		response = b.BaseBot.HandleCommand(message)
// 		/*response, err = handleCommand(sessionID, message, b)
// 		if err != nil {
// 			fmt.Printf("An error occurred: %s \n", err.Error())
// 			response = "An error occurred while processing your command."
// 		}*/
// 	} else if screaming && len(message) > 0 {
// 		response = strings.ToUpper(message)
// 	} else {
// 		// Get all document embeddings
// 		documentEmbeddings, chunkText, err := b.BaseBot.dao.FetchEmbeddings()
// 		//documentEmbeddings, chunkText, err := b.Service.GetAllDocumentEmbeddings()
// 		if err != nil {
// 			fmt.Printf("Error retrieving document embeddings: %v", err)
// 			response = "Error retrieving document embeddings."
// 		} else if useOpenAI {
// 			// Perform similarity matching with the user's message
// 			topChunks, err := document.RetrieveTopNChunks(message, documentEmbeddings, b.embConfig.NumTopChunks, chunkText, b.embConfig.ScoreThreshold) // Retrieve top 3 relevant chunks
// 			if err != nil {
// 				fmt.Printf("Error retrieving document chunks: %v", err)
// 				response = "Error retrieving related document information."
// 			} else if len(topChunks) > 0 {
// 				// If there are similar chunks found, provide them as context for GPT
// 				context := strings.Join(topChunks, "\n")
// 				gptPrompt := fmt.Sprintf("Context:\n%s\nUser query: %s", context, message)

// 				// Call GPT with the context and user query
// 				response, err = b.BaseBot.GetOpenAIResponse(gptPrompt)
// 				if err != nil {
// 					response = fmt.Sprintf("OpenAI Error: %v", err)
// 				}
// 			} else {
// 				// If no relevant document found, fallback to OpenAI response
// 				response, err = b.BaseBot.GetOpenAIResponse(message)
// 				if err != nil {
// 					response = fmt.Sprintf("OpenAI Error: %v", err)
// 				}
// 			}

// 		} else {
// 			//response = fmt.Sprintf("You said: %s", message)
// 			//HandleDialogflowIntent(message string) (string, error) {
// 			b.BaseBot.handleMessageDialogflow(GENERAL, sessionID, message, b)
// 		}
// 	}

// 	if response != "" {
// 		fmt.Printf("Sent message %s \n", response)
// 		err := b.sendResponse(sessionID, response)
// 		if err != nil {
// 			//c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
// 			fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
// 		}
// 	}

// }

func (b *generalBot) SendResponse(identifier interface{}, response string) error {
	// Perform type assertion to convert identifier to string
	if sessionID, ok := identifier.(string); ok {
		// Retrieve context using the sessionID
		c, err := getContext(sessionID)
		if err != nil {
			return fmt.Errorf("failed to retrieve context for sessionID: %s, error: %w", sessionID, err)
		}
		// Call sendFrontendMessage using the retrieved context
		return b.sendFrontendMessage(c, response)
	}
	return fmt.Errorf("invalid identifier type, expected string")
}

func (b *generalBot) sendFrontendMessage(c *gin.Context, message string) error {
	if c == nil {
		return fmt.Errorf("gin context is nil")
	}
	c.JSON(http.StatusOK, gin.H{
		"response": message,
	})
	return nil
}

var sessionContextMap = make(map[string]*gin.Context)

// StoreContext stores the context in sessionContextMap using the session ID
func (b *generalBot) StoreContext(sessionID string, c *gin.Context) {
	sessionContextMap[sessionID] = c
}

// Retrieve the context using sessionID when you need to send a response
func getContext(sessionID string) (*gin.Context, error) {
	if context, ok := sessionContextMap[sessionID]; ok {
		return context, nil
	}
	return nil, fmt.Errorf("no context found for session ID %s", sessionID)
}

func (b *generalBot) handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error {
	// Send the response to the respective platform or frontend
	for _, msg := range response.QueryResult.FulfillmentMessages {
		if text := msg.GetText(); text != nil {
			return b.SendResponse(identifier, text.Text[0])
		}
	}
	return fmt.Errorf("invalid identifier for frontend or platform")
}

/*func (b *generalBot) ProcessDocument(filename, sessionID, filePath string) ([]Document, []string, error) {
	// Extract text from the uploaded file (assuming downloadAndExtractText can handle local files)
	docText, err := document.DownloadAndExtractText(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing document: %w", err)
	}

	chunks := document.OverlapChunk(docText, b.embConfig.ChunkSize, b.embConfig.MinChunkSize)

	documents := make([]Document, 0)

	tagList := []string{} // Initialize the tag list

	for i, chunk := range chunks {

		document, err := b.StoreDocumentChunks(filename, filename+"_"+sessionID, chunk, i)
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

func (b *generalBot) StoreDocumentChunks(filename, docID, chunkText string, chunkid int) (Document, error) {
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

// func (b *generalBot) V2ProcessDocument(filename, sessionID, filePath string) ([]Document, []string, error) {
// 	// Extract text from the uploaded file (assuming downloadAndExtractText can handle local files)
// 	docText, err := document.DownloadAndExtractText(filePath)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("error processing document: %w", err)
// 	}

// 	documents, err := b.V2StoreDocumentChunks(filename, filename+"_"+sessionID, docText, b.embConfig.ChunkSize, b.embConfig.MinChunkSize)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Auto-tagging using OpenAI
// 	tags, err := b.openAIclient.AutoTagWithOpenAI(docText)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("error auto-tagging document: %w", err)
// 	}

// 	return documents, tags, nil
// }

// func (b *generalBot) V2StoreDocumentChunks(filename, docID, docText string, chunkSize, overlap int) ([]Document, error) {
// 	// Chunk the document with overlap
// 	chunks := document.OverlapChunk(docText, chunkSize, overlap)

// 	//client := openai.NewClient()

// 	documents := make([]Document, 0)
// 	for i, chunk := range chunks {
// 		// Get the embeddings for each chunk
// 		embedding, err := b.openAIclient.EmbedText(chunk)
// 		if err != nil {
// 			return nil, fmt.Errorf("error embedding chunk %d: %v", i, err)
// 		}

// 		// Create a unique chunk ID for storage in the database
// 		chunkID := fmt.Sprintf("%s_chunk_%d_%s", filename, i, docID)

// 		docText = utils.SanitizeText(docText)
// 		embeddingStr := utils.Float64SliceToPostgresArray(embedding)

// 		document := Document{
// 			Filename: filename,
// 			DocID:    docID,
// 			ChunkID:  chunkID,
// 			//DocText:   docText,
// 			DocText:   chunk,
// 			Embedding: embeddingStr,
// 		}
// 		documents = append(documents, document)
// 	}

// 	return documents, nil
// }
