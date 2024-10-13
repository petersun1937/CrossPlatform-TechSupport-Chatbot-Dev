package service

import (
	"crossplatform_chatbot/bot"
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	"crossplatform_chatbot/models"
	"crossplatform_chatbot/utils"
	"crossplatform_chatbot/utils/token"
	"fmt"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

type Service struct {
	database database.Database
	bots     map[string]bot.Bot
}

func NewService(botConfig config.BotConfig, database database.Database) *Service {
	return &Service{
		bots:     createBots(botConfig, database),
		database: database,
	}
}

func (s *Service) Init() error {
	// running bots
	for _, bot := range s.bots {
		if err := bot.Run(); err != nil {
			// log.Fatal("running bot failed:", err)
			fmt.Printf("running bot failed: %s", err.Error())
			return err
		}
	}
	return nil
}

func (s *Service) GetBot(tag string) bot.Bot {
	return s.bots[tag]
}

func createBots(botConfig config.BotConfig, database database.Database) map[string]bot.Bot {
	// Initialize bots
	lineBot, err := bot.V2NewLineBot(botConfig, database)
	if err != nil {
		//log.Fatal("Failed to initialize LINE bot:", err)
		fmt.Printf("Failed to initialize LINE bot: %s", err.Error())
	}

	// tgBot, err := bot.NewTGBot(conf, srv)
	// if err != nil {
	// 	//log.Fatal("Failed to initialize Telegram bot:", err)
	// 	fmt.Printf("Failed to initialize Telegram bot: %s", err.Error())
	// }

	// fbBot, err := bot.NewFBBot(conf, srv)
	// if err != nil {
	// 	log.Fatalf("Failed to create Facebook bot: %v", err)
	// }

	// igBot, err := bot.NewIGBot(conf, srv)
	// if err != nil {
	// 	log.Fatalf("Failed to create Instagram bot: %v", err)
	// }

	// generalBot, err := bot.NewGeneralBot(conf, srv)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize General bot: %v", err)
	// }

	return map[string]bot.Bot{
		"line": lineBot,
		// "tg":      tgBot,
		// "fb":      fbBot,
		// "ig":      igBot,
		// "general": generalBot,
	}
}

type ValidateUserReq struct {
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string
}

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        string `json:"id"` // Facebook User ID
}

// GetDB returns the gorm.DB instance from the service's database
func (s *Service) GetDB() *gorm.DB {
	return s.database.GetDB()
}

func (s *Service) ValidateUser(userIDStr string, req ValidateUserReq) (*string, error) {
	repo := NewRepository(s.database)

	// Check if the user exists
	_, err := repo.GetUser(userIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// If user not found, create a new user
			err = repo.CreateUser(userIDStr, req)
			if err != nil {
				return nil, err
			}

			// Generate a JWT token for the new user
			token, err := token.GenerateToken(userIDStr, "user") // Ensure GenerateToken accepts string
			if err != nil {
				fmt.Printf("Error generating JWT: %s", err.Error())
				return nil, err
			}

			return &token, nil
		}

		// Other errors when fetching user
		return nil, err
	}

	return nil, err
}

// StoreDocumentEmbedding stores the document and its embedding into the database
func (s *Service) StoreDocumentEmbedding(docID string, docText string, embedding []float64) error {
	repo := NewRepository(s.database)

	// Sanitize the document text
	docText = sanitizeText(docText)

	// Convert embedding to PostgreSQL-compatible array string
	embeddingStr := utils.Float64SliceToPostgresArray(embedding)

	docEmbedding := models.DocumentEmbedding{
		DocID:     docID,
		DocText:   docText,
		Embedding: embeddingStr,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the sanitized document embedding
	err := repo.CreateDocumentEmbedding(&docEmbedding)
	if err != nil {
		return fmt.Errorf("error storing document embedding: %v", err)
	}

	return nil
}

// CreateDocumentEmbedding stores a document embedding in the database
func (r *Repository) CreateDocumentEmbedding(docEmbedding *models.DocumentEmbedding) error {
	return r.database.GetDB().Create(docEmbedding).Error
}

func (r *Repository) BatchCreateDocumentEmbedding(docEmbeddings []*models.DocumentEmbedding) error {
	return r.database.GetDB().Create(docEmbeddings).Error
}

func sanitizeText(input string) string {
	validRunes := []rune{}
	for _, r := range input {
		if r == utf8.RuneError {
			continue // Skip invalid characters
		}
		validRunes = append(validRunes, r)
	}
	return string(validRunes)
}

// func (s *Service) ValidateUser(userIDStr string, req ValidateUserReq) (*string, error) {
// 	// Check if the user exists in the database
// 	var dbUser models.User
// 	// err := s.dao.CreatePlayer()
// 	err := s.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error

// 	// If the user does not exist, create a new user record
// 	if err != nil {

// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			// User does not exist, create a new user record
// 			dbUser = models.User{
// 				Model: gorm.Model{
// 					ID: 1,
// 				},
// 				UserID:       userIDStr,
// 				FirstName:    req.FirstName,
// 				LastName:     req.LastName,
// 				UserName:     req.UserName,
// 				LanguageCode: req.LanguageCode,
// 			}

// 			//err = database.DB.Create(&dbUser).Error
// 			err = s.database.GetDB().Create(&dbUser).Error

// 			if err != nil {
// 				fmt.Printf("Error creating user: %s", err.Error())
// 				return nil, err
// 			}

// 			// Generate a JWT token for the new user   TODO: move GenerateToken to bot?
// 			token, err := token.GenerateToken(userIDStr, "user") // Ensure GenerateToken accepts string
// 			if err != nil {
// 				fmt.Printf("Error generating JWT: %s", err.Error())
// 				return nil, err
// 			}

// 			return &token, nil

// 			// // Send the token to the user
// 			// msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Your access token is: "+token)
// 			// utils.TgBot.Send(msg)
// 		} else {
// 			// Handle other types of errors
// 			fmt.Printf("Error retrieving user: %s", err.Error())
// 		}

// 	}

// 	return nil, err
// }
