package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// TelegramUpdate represents a Telegram update
type TelegramUpdate struct {
	tgbotapi.Update
	/*UpdateID int             `json:"update_id"`
	Message  TelegramMessage `json:"message"`*/
}

// TelegramMessage represents a Telegram message
type TelegramMessage struct {
	*tgbotapi.Message
	/*MessageID int          `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
	Date      int          `json:"date"`
	Text      string       `json:"text"`*/
}

// TelegramUser represents a Telegram user
type TelegramUser struct {
	*tgbotapi.User
	/*
		ID           int64  `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		UserName     string `json:"username"`
		LanguageCode string `json:"language_code"`*/
}

// TelegramChat represents a Telegram chat
type TelegramChat struct {
	*tgbotapi.Chat
	/*
		ID        int64  `json:"id"`
		Type      string `json:"type"`
		Title     string `json:"title"`
		UserName  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`*/
}
