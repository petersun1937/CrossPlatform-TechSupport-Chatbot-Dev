package service

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"
	"errors"
	"fmt"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/line/line-bot-sdk-go/linebot"
)

// HandleLine processes incoming requests from the LINE platform.
func (s *Service) HandleLine(req *http.Request) error {
	lineBot, exist := s.GetBot("line").(bot.LineBot)
	if !exist {
		return errors.New("line bot not found")
	}

	// Parse the incoming request from the Line platform and extract the events
	events, err := lineBot.ParseRequest(req)
	if err != nil {
		return err
	}

	// Loop through each event received from the Line platform
	for _, event := range events {
		// Check if the event is a message event
		if event.Type == linebot.EventTypeMessage {
			// Switch on the type of message (could be text, image, video, audio etc., only support text for now)
			switch message := event.Message.(type) {
			// If the message is a text message, process it using handleLineMessage
			case *linebot.TextMessage:
				lineBot.HandleLineMessage(event, message)
			}
		}
	}

	return nil
}

// HandleTelegram processes incoming updates from Telegram, including documents and messages.
func (s *Service) HandleTelegram(update tgbotapi.Update) error {
	tgBot, exists := s.GetBot("telegram").(bot.TgBot)
	if !exists {
		return errors.New(" Telegram bot not found")
	}

	//tgBot.HandleTelegramUpdate(update)

	if update.Message != nil {
		if update.Message.Document != nil {

			// get filename, fileURL, fileID
			fileID, fileURL, filename, err := tgBot.GetDocFile(update)
			if err != nil {
				return fmt.Errorf("error getting file:  %w", err)
			}

			// If the message contains a document, handle the document upload
			err = s.HandleTGDocumentUpload(filename, fileID, fileURL)
			if err != nil {
				tgBot.SendTelegramMessage(update.Message.Chat.ID, "Error handling document: "+err.Error())
				return fmt.Errorf("error handling the document:  %w", err)
			}

			tgBot.SendTelegramMessage(update.Message.Chat.ID, "Document processed and stored in chunks for future queries.")
		} else {
			// Otherwise, handle regular text messages
			tgBot.HandleTgMessage(update.Message)
		}
	}

	return nil
}

// HandleMessenger processes incoming events from Facebook Messenger.
func (s *Service) HandleMessenger(event bot.MessengerEvent) error {
	fbBot, exists := s.GetBot("facebook").(bot.FbBot)
	if !exists {
		return errors.New(" Messenger bot not found")
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if messageText := strings.TrimSpace(msg.Message.Text); messageText != "" {
				fbBot.HandleMessengerMessage(senderID, messageText)
			}
		}
	}
	return nil
}

// HandleInstagram processes incoming events from Instagram.
func (s *Service) HandleInstagram(event bot.InstagramEvent) error {
	igBot, exists := s.GetBot("instagram").(bot.IgBot)
	if !exists {
		return errors.New(" Instagram bot not found")
	}

	for _, entry := range event.Entry {
		for _, msg := range entry.Messaging {
			senderID := msg.Sender.ID
			if messageText := strings.TrimSpace(msg.Message.Text); messageText != "" {
				igBot.HandleInstagramMessage(senderID, messageText)
			}
		}
	}
	return nil
}

// HandleGeneral processes requests from the frontend for the general bot.
func (s *Service) HandleGeneral(req models.GeneralRequest) error {

	// Get the bot responsible for sending the response.
	b := s.GetBot("general")
	/*_, exists := b.(bot.GeneralBot)
	if !exists {
		return errors.New("general bot not found")
	}*/

	// Process the message and generate a response using the service layer.
	response, err := s.ProcessUserMessage(req.SessionID, req.Message, "general")
	if err != nil {
		return fmt.Errorf("error processing user message: %w", err)
	}

	// Use the bot to send the response.
	//if response != "" {
	fmt.Printf("Sent message %s \n", response)
	err = b.SendResponse(req.SessionID, response)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while sending the response"})
		fmt.Printf("An error occurred while sending the response: %s\n", err.Error())
	}
	//}

	return nil
}

// func (s *Service) HandleGeneral(req models.GeneralRequest) error {

// 	genBot, exists := s.GetBot("general").(bot.GeneralBot)
// 	if !exists {
// 		return errors.New("general bot not found")
// 	}

// 	// Delegate the handling of the message to the general bot.
// 	genBot.HandleGeneralMessage(req.SessionID, req.Message)

// 	return nil
// }

func (s *Service) GetBotPlatform(botTag string) (bot.Bot, bot.Platform, error) {
	bot := s.GetBot(botTag)
	if bot == nil {
		return nil, 0, fmt.Errorf("bot not found for tag: %s", botTag)
	}
	return bot, bot.Platform(), nil
}
