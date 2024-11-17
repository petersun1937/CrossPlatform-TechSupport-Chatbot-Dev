package service

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) HandleDocumentUpload(filename, filePath string) error {
	// step 1: call bot to process documents
	b := s.GetBot("general").(bot.GeneralBot)

	// Generate a unique document ID
	uniqueDocID := uuid.New().String()

	documents, tags, err := b.V2ProcessDocument(filename, uniqueDocID, filePath)
	if err != nil {
		return err
	}

	// dao version
	// return s.repository.CreateDocumentsAndMeta(uniqueDocID, documents, tags)

	// service version
	// step 2: make db data
	documentModels := make([]*models.Document, 0)
	for _, doc := range documents {
		model := models.Document{
			Filename:  doc.Filename,
			DocID:     doc.DocID,
			ChunkID:   doc.ChunkID,
			DocText:   doc.DocText,
			Embedding: doc.Embedding,
		}
		documentModels = append(documentModels, &model)
	}
	metadata := models.DocumentMetadata{
		DocID: uniqueDocID,
		Tags:  tags,
	}

	// step 3: do transaction
	return s.database.GetDB().Transaction(func(tx *gorm.DB) error {
		// batch insert Documents
		if err := tx.Create(documentModels).Error; err != nil {
			return err
		}

		// insert DocumentMetadata
		if err := tx.Create(&metadata).Error; err != nil {
			return err
		}

		return nil
	})
}
