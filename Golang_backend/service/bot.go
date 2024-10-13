package service

import (
	"errors"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

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
