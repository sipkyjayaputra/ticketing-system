package usecase

import (
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"
)

func (uc *usecase) GetTickets() (*utils.ResponseContainer, *utils.ErrorContainer) {
	// tickets, err := uc.repo.GetTickets()

	// if err != nil {
	// 	return nil, utils.BuildInternalErrorResponse("failed to get tickets", err.Error())
	// }

	// Optionally, you can include any extra processing needed for activities and documents here.

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) AddTicket(ticket dto.Ticket, creator uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	ticket.CreatedBy = creator
	ticket.UpdatedBy = creator

	// Handle activities and documents within the ticket if necessary
	for i := range ticket.Activities {
		ticket.Activities[i].CreatedAt = time.Now()
		ticket.Activities[i].UpdatedAt = time.Now()
		ticket.Activities[i].CreatedBy = creator
		ticket.Activities[i].UpdatedBy = creator

		// For each activity, set the creator for its documents
		for j := range ticket.Activities[i].Files {
			ticket.Activities[i].Files[j].CreatedAt = time.Now()
			ticket.Activities[i].Files[j].UpdatedAt = time.Now()
			ticket.Activities[i].Files[j].CreatedBy = creator
			ticket.Activities[i].Files[j].UpdatedBy = creator
		}
	}

	if err := uc.repo.AddTicket(ticket); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to add ticket", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) UpdateTicket(ticket dto.Ticket, updater uint, ticketNo uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	ticket.UpdatedBy = updater

	// Handle activities and documents within the ticket if necessary
	for i := range ticket.Activities {
		ticket.Activities[i].UpdatedBy = updater

		// For each activity, update the updater for its documents
		for j := range ticket.Activities[i].Files {
			ticket.Activities[i].Files[j].UpdatedBy = updater
		}
	}

	// if err := uc.repo.UpdateTicket(ticket, ticketNo); err != nil {
	// 	return nil, utils.BuildInternalErrorResponse("failed to update ticket", err.Error())
	// }

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) DeleteTicket(ticketNo uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	// if err := uc.repo.DeleteTicket(ticketNo); err != nil {
	// 	return nil, utils.BuildInternalErrorResponse("failed to delete ticket", err.Error())
	// }

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) GetTicketById(ticketNo uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	// ticket, err := uc.repo.GetTicketById(ticketNo)

	// if err != nil {
	// 	return nil, utils.BuildInternalErrorResponse("failed to get ticket", err.Error())
	// }

	// Optionally, include any extra processing needed for activities and documents here.

	return utils.BuildSuccessResponse(nil), nil
}
