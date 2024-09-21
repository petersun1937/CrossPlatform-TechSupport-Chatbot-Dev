package bot

import (
	config "Tg_chatbot/configs"
	"Tg_chatbot/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type FbBot interface {
	HandleMessengerMessage(senderID, messageText string)
	Run() error
}

type fbBot struct {
	BaseBot
	ctx             context.Context
	pageAccessToken string
}

// creates a new FbBot instance
func NewFBBot(conf *config.Config, service service.Service) (*fbBot, error) {
	// Verify that the page access token is available
	if conf.FacebookPageToken == "" {
		return nil, errors.New("facebook Page Access Token is not provided")
	}

	// Initialize the BaseBot structure
	baseBot := &BaseBot{
		Platform: FACEBOOK,
		Service:  service,
	}

	// Initialize and return the fbBot instance
	return &fbBot{
		BaseBot:         baseBot,
		ctx:             context.Background(),
		pageAccessToken: conf.FacebookPageToken,
	}, nil
}

// Run initializes and starts the Facebook bot with webhook
func (b *fbBot) Run() error {
	if b.pageAccessToken == "" {
		return errors.New("page access token is missing")
	}

	//TODO

	// webhook confirmation
	fmt.Println("Facebook Messenger bot is running with webhook!")

	return nil
}

// MessengerEvent defines the structure of incoming events from Facebook Messenger
type MessengerEvent struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

// HandleMessengerMessage processes incoming messages and sends a response
func (b *fbBot) HandleMessengerMessage(senderID, messageText string) {
	responseText := messageText // TODO response logic

	// Call the sendMessage function to reply to the user
	err := b.sendMessengerResponse(senderID, responseText)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// sendMessage sends a message to the specified user on Messenger
func (b *fbBot) sendMessengerResponse(recipientID, messageText string) error {
	conf := config.GetConfig()
	url := conf.FacebookAPIURL + "/messages?access_token=" + b.pageAccessToken

	// Create the message payload
	messageData := map[string]interface{}{
		"recipient": map[string]string{"id": recipientID},
		"message":   map[string]string{"text": messageText},
	}

	// Marshal the payload to JSON
	messageBody, err := json.Marshal(messageData)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	// Create HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending response: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Message sent successfully to %s", recipientID)
	return nil
}
