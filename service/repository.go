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

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(userIDStr string, req ValidateUserReq) error {
	//err = database.DB.Create(&dbUser).Error
	return r.database.GetDB().Create(&models.User{
		Model:        gorm.Model{},
		UserID:       userIDStr,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		UserName:     req.UserName,
		LanguageCode: req.LanguageCode,
	}).Error
}

func (r *Repository) GetUser(userIDStr string) (*models.User, error) {
	var dbUser models.User
	err := r.database.GetDB().Where("user_id = ? AND deleted_at IS NULL", userIDStr).First(&dbUser).Error
	if err != nil {
		return nil, err
	}
	return &dbUser, nil
}
