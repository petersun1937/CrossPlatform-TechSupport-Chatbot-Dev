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
	repo := NewRepository(s.database)
	if err := repo.GetUser(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := repo.CreateUser(userID, userName); err != nil {
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
	}

	return nil, err
}
