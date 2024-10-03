package repository

import "github.com/sipkyjayaputra/ticketing-system/model/entity"

// GetDocumentById retrieves a specific document by ID from the database
func (repo *repository) GetDocumentById(documentID uint) (*entity.Document, error) {
	document := &entity.Document{}
	if err := repo.db.Model(&entity.Document{}).Where("id = ?", documentID).First(&document).Error; err != nil {
		return nil, err
	}
	return document, nil
}
