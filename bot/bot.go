package bot

import (
	"Tg_chatbot/service"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

type Bot interface {
	Run() error

	sendMenu(identifier interface{}) error
	sendResponse(identifier interface{}, message string) error
	handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error
}

type BaseBot struct {
	Platform Platform
	Service  *service.Service
}

// define platforms
type Platform int

const (
	LINE Platform = iota
	TELEGRAM
)
