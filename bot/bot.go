package bot

import (
	config "crossplatform_chatbot/configs"
	"crossplatform_chatbot/database"
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/repository"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// action, public function
type Bot interface {
	Run() error

	//sendMenu(identifier interface{}) error
	SendResponse(identifier interface{}, message string) error
	handleDialogflowResponse(response *dialogflowpb.DetectIntentResponse, identifier interface{}) error
	Base() *BaseBot // Make BaseBot accessible
	Platform() Platform
}

type BaseBot struct {
	platform Platform
	// Service  *service.Service
	conf         *config.BotConfig
	database     database.Database
	dao          repository.DAO
	openAIclient *openai.Client
	embConfig    config.EmbeddingConfig
}

func (b *BaseBot) Base() *BaseBot {
	return b
}

func (b *BaseBot) Platform() Platform {
	return b.platform
}

// define platforms
type Platform int

const (
	LINE Platform = iota
	TELEGRAM
	FACEBOOK
	INSTAGRAM
	GENERAL
)
