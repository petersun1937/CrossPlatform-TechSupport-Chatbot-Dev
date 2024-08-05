package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID       int64  `gorm:"unique"`
	FirstName    string `json:"firstname" bson:"firstname"`
	LastName     string `json:"lastname" bson:"lastname"`
	UserName     string `json:"username" bson:"username" binding:"required"`
	LanguageCode string
}
