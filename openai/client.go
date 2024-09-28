package openai

import (
	config "crossplatform_chatbot/configs"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// Struct to hold OpenAI API client configuration
type Client struct {
	ApiKey string
	Client *resty.Client
}

// Function to create a new OpenAI client
func NewClient() *Client {
	conf := config.GetConfig()
	apiKey := conf.OpenaiAPIKey
	if apiKey == "" {
		log.Fatal("OpenAI API key not found in environment variables")
	}

	client := resty.New()
	return &Client{
		ApiKey: apiKey,
		Client: client,
	}
}

// Function to get a response from the OpenAI API
func (c *Client) GetResponse(prompt string) (string, error) {
	request := map[string]interface{}{
		"model":       "gpt-3.5-turbo",                                          // Specify model type here (can't use gpt-4?)
		"messages":    []map[string]string{{"role": "user", "content": prompt}}, // Adjusted for chat models
		"max_tokens":  100,
		"temperature": 0.7,
	}

	// Send the request to OpenAI API (chat completion endpoint)
	response, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.openai.com/v1/chat/completions") // Ensure the correct endpoint

	if err != nil {
		return "", fmt.Errorf("error sending request to OpenAI: %v", err)
	}

	if response.StatusCode() != 200 {
		return "", fmt.Errorf("OpenAI API returned status code %d: %s", response.StatusCode(), response.String())
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return "", fmt.Errorf("error parsing response from OpenAI: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no response from OpenAI")
	}

	text, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", errors.New("invalid response format from OpenAI")
	}

	return text, nil
}
