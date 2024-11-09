package document

import (
	"math"
	"strings"

	levenshtein "github.com/texttheater/golang-levenshtein/levenshtein"
)

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


