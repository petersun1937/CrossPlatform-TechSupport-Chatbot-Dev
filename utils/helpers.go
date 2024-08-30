package utils

import (
	"context"
	"fmt"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
)

// Global variables to hold the bot instances for TG and LINE
//var TgBot *tgbotapi.BotAPI
//var LineBot *linebot.Client

// Send a text query to Dialogflow and returns the response
func DetectIntentText(projectID, sessionID, text, languageCode string) (*dialogflowpb.DetectIntentResponse, error) {
	// Create a background context for the API call
	ctx := context.Background()

	// Create a new Dialogflow Sessions client
	client, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close() // Ensure the client is closed when done

	// Construct the session path for the Dialogflow API
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)

	// Create the DetectIntentRequest with the session path and query input
	req := &dialogflowpb.DetectIntentRequest{
		Session: sessionPath,
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         text,
					LanguageCode: languageCode,
				},
			},
		},
	}

	// Send the request and return the response or error
	return client.DetectIntent(ctx, req)
}
