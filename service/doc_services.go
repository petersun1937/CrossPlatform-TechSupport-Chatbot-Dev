package service

import (
	"crossplatform_chatbot/bot"
	"crossplatform_chatbot/models"

	"gorm.io/gorm"
)

func (s *Service) GetUploadedDocuments() ([]string, error) {
	documents, err := s.repository.GetAllDocuments() // Fetch all documents from the repository
	if err != nil {
		return nil, err
	}

	// Use a map to store unique filenames
	uniqueFilenameMap := make(map[string]struct{})
	//uniqueFilenames := []string{}
	uniqueFilenames := make([]string, 0, len(documents))

	// Loop through documents and collect unique filenames
	for _, doc := range documents {
		filename := doc.Filename

		// Check if the filename is a duplicate
		/*if count, exists := filenameCount[filename]; exists {
			// Increment the counter and append it to the filename to make it unique
			newFilename := fmt.Sprintf("%s(%d)", filename, count+1)
			uniqueFilenames = append(uniqueFilenames, newFilename)
			filenameCount[filename] = count + 1
		} else {
			// If it's a unique filename, store it directly
			uniqueFilenames = append(uniqueFilenames, filename)
			filenameCount[filename] = 0
		}*/

		// If filename is not already in the map, add it to the list
		if _, exists := uniqueFilenameMap[filename]; !exists {
			uniqueFilenames = append(uniqueFilenames, filename)
			uniqueFilenameMap[filename] = struct{}{} // Store it in the map
		}
	}

	return uniqueFilenames, nil
}

func (s *Service) HandleDocumentUpload(filename, fileID, filePath string) error {
	// step 1: call bot to process documents
	b := s.GetBot("general").(bot.GeneralBot)

	documents, tags, err := b.ProcessDocument(filename, fileID, filePath)
	if err != nil {
		return err
	}

	// dao version
	// return s.repository.CreateDocumentsAndMeta(fileID, documents, tags)

	// service version
	// step 2: make db data
	documentModels := make([]*models.Document, 0)
	documentID := ""
	for _, doc := range documents {
		model := models.Document{
			Filename:  doc.Filename,
			DocID:     doc.DocID,
			ChunkID:   doc.ChunkID,
			DocText:   doc.DocText,
			Embedding: doc.Embedding,
		}
		documentModels = append(documentModels, &model)
		documentID = doc.DocID
	}
	metadata := models.DocumentMetadata{
		DocID: documentID,
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

func (s *Service) HandleTGDocumentUpload(filename, fileID, filePath string) error {
	// step 1: call bot to process documents
	b := s.GetBot("telegram").(bot.TgBot)

	documents, tags, err := b.ProcessDocument(filename, fileID, filePath)
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
		DocID: fileID,
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
