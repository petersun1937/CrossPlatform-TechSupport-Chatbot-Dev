package service

import (
	"Tg_chatbot/database"
	"Tg_chatbot/models"

	"gorm.io/gorm"
)

type Repository struct {
	database database.Database
}

func NewRepository(database database.Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) CreateUser() error {
	//err = database.DB.Create(&dbUser).Error
	return r.database.GetDB().Create(&models.User{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:       userIDStr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UserName:     req.UserName,
		LanguageCode: req.LanguageCode,
	}).Error
}

func (r *Repository) GetUser() (*models.User, error) {
	return s.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
}
