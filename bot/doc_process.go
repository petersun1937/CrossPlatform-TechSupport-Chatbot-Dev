package bot

import (
	openai "crossplatform_chatbot/openai"
	"crossplatform_chatbot/utils"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func downloadAndExtractText(fileURL string) (string, error) {
	response, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("error downloading file: %v", err)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading file content: %v", err)
	}

	// Assume the content is plain text for simplicity. You can add PDF or DOCX parsing logic here.
	return string(content), nil
}

//var documentEmbeddings = make(map[string][]float64)

func ChunkDocument(text string, chunkSize int) []string {
	words := strings.Fields(text) // Split the document into words
	var chunks []string
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	fmt.Printf("Document chunked into %d chunks.\n", len(chunks)) // Log the chunk count
	return chunks
}

func (b *tgBot) StoreDocumentChunks(docID string, text string, chunkSize int, minchunkSize int) error {
	//chunks := ChunkDocument(text, chunkSize)
	//chunks := utils.ChunkDocumentBySentence(text, chunkSize)
	chunks := utils.ChunkSmartly(text, chunkSize, minchunkSize)
	for i, chunk := range chunks {
		embedding, err := openai.EmbedDocument(chunk)
		if err != nil {
			return fmt.Errorf("error embedding chunk %d: %v", i, err)
		}
		chunkID := fmt.Sprintf("%s_chunk_%d", docID, i)
		b.Service.StoreDocumentEmbedding(chunkID, chunk, embedding) // Store each chunk with its embedding
	}
	fmt.Println("Document embedding complete.")
	return nil
}
