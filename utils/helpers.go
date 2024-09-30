package utils

import (
	"context"
	"crossplatform_chatbot/openai"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	levenshtein "github.com/texttheater/golang-levenshtein/levenshtein"
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

// Convert float64 slice to PostgreSQL float8[] string format
func Float64SliceToPostgresArray(embedding []float64) string {
	var result strings.Builder
	result.WriteString("{")
	for i, value := range embedding {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString(fmt.Sprintf("%f", value))
	}
	result.WriteString("}")
	return result.String()
}

// Convert data type to store embeddings in Postgres
func PostgresArrayToFloat64Slice(embeddingStr string) ([]float64, error) {
	// Remove curly braces from the string
	embeddingStr = strings.Trim(embeddingStr, "{}")

	// Split the string by commas
	stringValues := strings.Split(embeddingStr, ",")

	// Convert the string values back to float64
	floatValues := make([]float64, len(stringValues))
	for i, strVal := range stringValues {
		val, err := strconv.ParseFloat(strings.TrimSpace(strVal), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing embedding value: %v", err)
		}
		floatValues[i] = val
	}

	return floatValues, nil
}

// Compute similarity score and retrieve  the top N chunks from database.
func RetrieveTopNChunks(query string, documentEmbeddings map[string][]float64, topN int, docIDToText map[string]string, threshold float64) ([]string, error) {
	fmt.Println("Embedding query for similarity search...")
	queryEmbedding, err := openai.EmbedDocument(query)
	if err != nil {
		return nil, fmt.Errorf("error embedding query: %v", err)
	}

	fmt.Println("Calculating similarity between query and document chunks...")
	type chunkScore struct {
		chunkID string
		score   float64
	}
	var scores []chunkScore

	// Calculate similarity for each document chunk
	for chunkID, embedding := range documentEmbeddings {
		//score := cosineSimilarity(queryEmbedding, embedding)
		cosineScore := cosineSimilarity(queryEmbedding, embedding)
		keywordScore := keywordMatchScore(query, docIDToText[chunkID])
		combinedScore := weightedScore(cosineScore, keywordScore)

		fmt.Printf("Combined score for chunk %s: %f\n", chunkID, combinedScore)

		// Only add the chunk if it meets the threshold
		if combinedScore >= threshold {
			scores = append(scores, chunkScore{chunkID, combinedScore})
		}
		/*if score >= similarityThreshold { // Filter out low similarity scores
			scores = append(scores, chunkScore{chunkID, score})
			fmt.Printf("Similarity score for chunk %s: %f\n", chunkID, score)
		}*/
	}

	// Sort the chunks based on score (highest score first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Collect the top N chunks' actual text using docIDToText
	var topChunksText []string
	for i := 0; i < topN && i < len(scores); i++ {
		chunkID := scores[i].chunkID
		if text, exists := docIDToText[chunkID]; exists {
			topChunksText = append(topChunksText, text)
		} else {
			topChunksText = append(topChunksText, fmt.Sprintf("Text not found for chunk: %s", chunkID))
		}
	}

	fmt.Println("Top relevant chunks selected.")
	return topChunksText, nil
}

// Fuzzy match score between two words using Levenshtein distance
func fuzzyMatchScore(queryWord string, chunkWord string) float64 {
	// Compute Levenshtein distance between the query word and chunk word
	distance := levenshtein.DistanceForStrings([]rune(queryWord), []rune(chunkWord), levenshtein.DefaultOptions)

	// Calculate similarity ratio (1 - normalized distance)
	maxLen := math.Max(float64(len(queryWord)), float64(len(chunkWord)))
	if maxLen == 0 {
		return 0
	}
	return 1 - float64(distance)/maxLen
}

// Modified keywordMatchScore function with fuzzy matching
func keywordMatchScore(query string, chunkText string) float64 {
	queryWords := strings.Fields(strings.ToLower(query))
	chunkWords := strings.Fields(strings.ToLower(chunkText))

	matchCount := 0
	for _, queryWord := range queryWords {
		for _, chunkWord := range chunkWords {
			// Use a fuzzy match score with a threshold (0.8 for close matches)
			if fuzzyMatchScore(queryWord, chunkWord) > 0.8 {
				matchCount++
				break
			}
		}
	}

	return float64(matchCount) / float64(len(queryWords)) // Return a ratio of matching words
}

// Weighted score combined with Cosine similarity and fuzzy keyword matching
func weightedScore(cosineScore float64, keywordScore float64) float64 {
	return (0.7 * cosineScore) + (0.3 * keywordScore) // Weight cosine higher but consider keyword match
}

// Compute similarity score
func cosineSimilarity(vec1, vec2 []float64) float64 {
	var dotProduct, normA, normB float64
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i] // Calculate the dot product of vec1 and vec2
		normA += vec1[i] * vec1[i]      // Calculate the sum of squares of vec1
		normB += vec2[i] * vec2[i]      // Calculate the sum of squares of vec2
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)) // Return the cosine similarity
}

// Chunks the text by full sentences, keeping each chunk under a certain word limit
func ChunkDocumentBySentence(text string, chunkSize int) []string {
	sentences := splitIntoSentences(text) // Split the document into sentences
	var chunks []string
	var currentChunk []string
	currentWordCount := 0

	// Iterate over the sentences and add them to chunks
	for _, sentence := range sentences {
		wordCount := len(strings.Fields(sentence)) // Count the words in the sentence

		// If adding this sentence exceeds the chunk size, start a new chunk
		if currentWordCount+wordCount > chunkSize && len(currentChunk) > 0 {
			chunks = append(chunks, strings.Join(currentChunk, " "))
			currentChunk = []string{}
			currentWordCount = 0
		}

		currentChunk = append(currentChunk, sentence)
		currentWordCount += wordCount
	}

	// Add the last chunk if it's not empty
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}

// Helper function to split text into paragraphs
func splitIntoParagraphs(text string) []string {
	// Split text into paragraphs using two or more newlines as the delimiter
	return strings.Split(text, "\n\n")
}

// Helper function to check if a string is an enumerated point (e.g., "1. ", "2. ", etc.)
func isEnumeratedPoint(line string) bool {
	re := regexp.MustCompile(`^\d+\.\s`)
	return re.MatchString(line)
}

// Helper function to chunk text into paragraphs and points
func ChunkSmartly(text string, maxChunkSize int, minWordsPerChunk int) []string {
	var chunks []string

	// Split the text into paragraphs first
	paragraphs := splitIntoParagraphs(text)
	for _, paragraph := range paragraphs {
		// Split paragraph into lines to detect enumerated points
		lines := strings.Split(paragraph, "\n")
		pointGroup := ""
		for _, line := range lines {
			line = strings.TrimSpace(line) // Ensure that each line is trimmed
			if line == "" {
				continue // Skip empty lines
			}

			// If the line starts with an enumerated point, group it
			if isEnumeratedPoint(line) {
				if pointGroup != "" {
					// Add the previous point group to chunks
					chunks = append(chunks, strings.TrimSpace(pointGroup))
					pointGroup = ""
				}
				pointGroup += line + " " // Group lines with enumerated points
			} else {
				// If a non-enumerated point, check if it's a new paragraph or continuation
				if pointGroup != "" {
					pointGroup += line + " "
				} else {
					// Add standalone sentences or paragraphs directly
					if len(line) <= maxChunkSize {
						chunks = append(chunks, strings.TrimSpace(line))
					} else {
						// Split into smaller sentences if it's still too long
						chunks = append(chunks, splitIntoSentences(line)...)
					}
				}
			}
		}
		if pointGroup != "" {
			chunks = append(chunks, strings.TrimSpace(pointGroup)) // Add the last grouped point
		}
	}

	// Combine chunks that are smaller than the minimum word count
	var combinedChunks []string
	currentChunk := ""
	currentWordCount := 0

	for _, chunk := range chunks {
		wordsInChunk := len(strings.Fields(chunk))

		// If adding this chunk doesn't exceed the minimum word count, combine it with the current chunk
		if currentWordCount+wordsInChunk < minWordsPerChunk {
			currentChunk += " " + chunk
			currentWordCount += wordsInChunk
		} else {
			// If current chunk meets the minimum size, append it to combinedChunks
			if len(strings.TrimSpace(currentChunk)) > 0 {
				combinedChunks = append(combinedChunks, strings.TrimSpace(currentChunk))
			}
			// Start a new chunk
			currentChunk = chunk
			currentWordCount = wordsInChunk
		}
	}

	// Add the last chunk if it hasn't been added yet
	if len(strings.TrimSpace(currentChunk)) > 0 {
		combinedChunks = append(combinedChunks, strings.TrimSpace(currentChunk))
	}

	return combinedChunks
}

// Helper function to split text into sentences, accounting for decimals
func splitIntoSentences(text string) []string {
	var sentences []string
	sentence := ""
	for i, r := range text {
		sentence += string(r)

		// Check for sentence-ending punctuation ('.', '!', '?')
		if r == '.' || r == '!' || r == '?' {
			// Ensure it's not a decimal point by checking if the character before the period is a digit
			if r == '.' && i > 0 && unicode.IsDigit(rune(text[i-1])) {
				continue // Skip splitting here if it's part of a number (e.g., 1.1, 1.2)
			}

			// Add the sentence to the list and reset the sentence accumulator
			sentences = append(sentences, strings.TrimSpace(sentence))
			sentence = ""
		}
	}

	// Add any remaining text as the last sentence
	if len(strings.TrimSpace(sentence)) > 0 {
		sentences = append(sentences, strings.TrimSpace(sentence))
	}

	return sentences
}
