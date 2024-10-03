package usecase

import (
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"github.com/sipkyjayaputra/ticketing-system/utils"
)

// GetDocumentById retrieves a specific activity by its ID
func (uc *usecase) GetDocumentById(documentID uint) (*entity.Document, *utils.ErrorContainer) {
	document, err := uc.repo.GetDocumentById(documentID)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get activity", err.Error())
	}

	return document, nil
}
