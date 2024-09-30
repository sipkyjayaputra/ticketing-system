package usecase

import (
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/repository"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/sirupsen/logrus"
)

type Usecase interface {
	// AUTH
	RefreshToken(string) (*utils.ResponseContainer, *utils.ErrorContainer)
	SignIn(dto.SignIn) (*utils.ResponseContainer, *utils.ErrorContainer)

	// USER
	GetUsers() (*utils.ResponseContainer, *utils.ErrorContainer)
	AddUser(dto.User) (*utils.ResponseContainer, *utils.ErrorContainer)
	UpdateUser(dto.User, string, string) (*utils.ResponseContainer, *utils.ErrorContainer)
	DeleteUser(string) (*utils.ResponseContainer, *utils.ErrorContainer)
	GetUserById(string) (*utils.ResponseContainer, *utils.ErrorContainer)

	// TICKET
	GetTickets(dto.TicketFilter) (*utils.ResponseContainer, *utils.ErrorContainer)
	AddTicket(dto.Ticket, uint) (*utils.ResponseContainer, *utils.ErrorContainer)
	UpdateTicket(dto.Ticket, uint, string) (*utils.ResponseContainer, *utils.ErrorContainer)
	DeleteTicket(string) (*utils.ResponseContainer, *utils.ErrorContainer)
	GetTicketById(string) (*utils.ResponseContainer, *utils.ErrorContainer)

	// ACTIVITY
	GetActivitiesByTicketNo(string) (*utils.ResponseContainer, *utils.ErrorContainer)
	AddActivity(dto.Activity, uint) (*utils.ResponseContainer, *utils.ErrorContainer)
	UpdateActivity(dto.Activity, uint, uint) (*utils.ResponseContainer, *utils.ErrorContainer)
	DeleteActivity(uint) (*utils.ResponseContainer, *utils.ErrorContainer)
	GetActivityById(uint) (*utils.ResponseContainer, *utils.ErrorContainer)
}

type usecase struct {
	repo   repository.Repository
	logger *logrus.Logger
}

func NewUsecase(repo repository.Repository, logger *logrus.Logger) Usecase {
	return &usecase{
		repo:   repo,
		logger: logger,
	}
}
