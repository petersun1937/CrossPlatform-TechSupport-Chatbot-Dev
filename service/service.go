package service

import (
	"Tg_chatbot/database"
	"Tg_chatbot/models"
	"Tg_chatbot/utils/token"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	database database.Database
}

func NewService(database database.Database) *Service {
	return &Service{
		database: database,
	}
}

type ValidateUserReq struct {
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string
}

// GetDB returns the gorm.DB instance from the service's database
func (s *Service) GetDB() *gorm.DB {
	return s.database.GetDB()
}

func (s *Service) ValidateUser(userIDStr string, req ValidateUserReq) (*string, error) {
	// Check if the user exists in the database
	var dbUser models.User
	// err := s.dao.CreatePlayer()
	err := s.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error

	// If the user does not exist, create a new user record
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist, create a new user record
			dbUser = models.User{
				Model: gorm.Model{
					ID: 1,
				},
				UserID:       userIDStr,
				FirstName:    req.FirstName,
				LastName:     req.LastName,
				UserName:     req.UserName,
				LanguageCode: req.LanguageCode,
			}

			//err = database.DB.Create(&dbUser).Error
			err = s.database.GetDB().Create(&dbUser).Error

			if err != nil {
				fmt.Printf("Error creating user: %s", err.Error())
				return nil, err
			}

			// Generate a JWT token for the new user
			token, err := token.GenerateToken(userIDStr, "user") // Ensure GenerateToken accepts string
			if err != nil {
				fmt.Printf("Error generating JWT: %s", err.Error())
				return nil, err
			}

			return &token, nil

			// // Send the token to the user
			// msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Your access token is: "+token)
			// utils.TgBot.Send(msg)
		} else {
			// Handle other types of errors
			fmt.Printf("Error retrieving user: %s", err.Error())
		}

	}

	return nil, err
}
