package models

// TelegramUpdate represents a Telegram update
type TelegramUpdate struct {
	UpdateID int             `json:"update_id"`
	Message  TelegramMessage `json:"message"`
}

// TelegramMessage represents a Telegram message
type TelegramMessage struct {
	MessageID int          `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
	Date      int          `json:"date"`
	Text      string       `json:"text"`
}

// TelegramUser represents a Telegram user
type TelegramUser struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// TelegramChat represents a Telegram chat
type TelegramChat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
