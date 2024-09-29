package usecase

import (
	"errors"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"
)

func (uc *usecase) GetTickets() (*utils.ResponseContainer, *utils.ErrorContainer) {
	tickets, err := uc.repo.GetTickets()

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get tickets", err.Error())
	}

	for i := range tickets {
		tickets[i].Assigned.Password = ""
		tickets[i].Assigned.CreatedAt = nil
		tickets[i].Assigned.UpdatedAt = nil
		tickets[i].Assigned.CreatedBy = ""
		tickets[i].Assigned.UpdatedBy = ""
		tickets[i].Reporter.Password = ""
		tickets[i].Reporter.CreatedAt = nil
		tickets[i].Reporter.UpdatedAt = nil
		tickets[i].Reporter.CreatedBy = ""
		tickets[i].Reporter.UpdatedBy = ""
	}

	return utils.BuildSuccessResponse(tickets), nil
}

func (uc *usecase) AddTicket(ticket dto.Ticket, creator uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	ticket.CreatedBy = creator
	ticket.UpdatedBy = creator
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	if len(ticket.Activities) < 1 {
		err := errors.New("no activity")
		return nil, utils.BuildInternalErrorResponse("failed to add ticket", err.Error())
	}

	if err := uc.repo.AddTicket(ticket); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to add ticket", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) UpdateTicket(ticket dto.Ticket, updater uint, ticketNo string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	ticket.TicketNo = ticketNo
	ticket.UpdatedBy = updater
	ticket.UpdatedAt = time.Now()

	if err := uc.repo.UpdateTicket(ticket); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to update ticket", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) DeleteTicket(ticketNo string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	if err := uc.repo.DeleteTicket(ticketNo); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to delete ticket", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) GetTicketById(ticketNo string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	ticket, err := uc.repo.GetTicketById(ticketNo)

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get ticket", err.Error())
	}

	ticket.Assigned.Password = ""
	ticket.Assigned.CreatedAt = nil
	ticket.Assigned.UpdatedAt = nil
	ticket.Assigned.CreatedBy = ""
	ticket.Assigned.UpdatedBy = ""
	ticket.Reporter.Password = ""
	ticket.Reporter.CreatedAt = nil
	ticket.Reporter.UpdatedAt = nil
	ticket.Reporter.CreatedBy = ""
	ticket.Reporter.UpdatedBy = ""

	return utils.BuildSuccessResponse(ticket), nil
}
