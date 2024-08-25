package handlers

/*
// mock of the Telegram bot API
type MockTelegramBot struct {
	mock.Mock
}

// mock implementation for testing
var mockTgBot *MockTelegramBot

func init() {
	mockTgBot = new(MockTelegramBot)
	utils.TgBot = mockTgBot
}

func TestHandleButton(t *testing.T) {
	mockBot := new(MockTelegramBot)

	query := &tgbotapi.CallbackQuery{
		ID:      "callback_id",
		Data:    utils.NextButton,
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 123}},
	}

	mockBot.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)

	handleButton(query)

	// Check if the correct response was sent
	mockBot.AssertExpectations(t)
}*/
