package openai

import (
	"encoding/json"
	"fmt"
)

// EmbedDocument converts text to an embedding vector using OpenAI's embedding model
func EmbedDocument(text string) ([]float64, error) {
	//fmt.Println("Starting document embedding process...")

	client := NewClient() // Create a new OpenAI client
	request := map[string]interface{}{
		"model": "text-embedding-ada-002",
		"input": text,
	}

	response, err := client.Client.R().
		SetHeader("Authorization", "Bearer "+client.ApiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post("https://api.openai.com/v1/embeddings")

	if err != nil {
		return nil, fmt.Errorf("error embedding document: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Body(), &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	data := result["data"].([]interface{})
	embedding := data[0].(map[string]interface{})["embedding"].([]interface{})

	//fmt.Println("Converting embedding to []float64...")
	embeddingFloat := make([]float64, len(embedding))
	for i, v := range embedding {
		embeddingFloat[i] = v.(float64)
	}

	//fmt.Println("Document embedding complete.")
	return embeddingFloat, nil
}

// func EmbedDocument(text string) ([]float64, error) {
// 	client := NewClient() // Create a new OpenAI client
// 	request := map[string]interface{}{
// 		"model": "text-embedding-ada-002",
// 		"input": text,
// 	}

// 	response, err := client.Client.R().
// 		SetHeader("Authorization", "Bearer "+client.ApiKey).
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(request).
// 		Post("https://api.openai.com/v1/embeddings")

// 	if err != nil {
// 		return nil, fmt.Errorf("error embedding document: %v", err)
// 	}

// 	var result map[string]interface{}
// 	if err := json.Unmarshal(response.Body(), &result); err != nil {
// 		return nil, fmt.Errorf("error parsing response: %v", err)
// 	}

// 	embedding := result["data"].([]interface{})[0].(map[string]interface{})["embedding"].([]float64)
// 	return embedding, nil
// }
