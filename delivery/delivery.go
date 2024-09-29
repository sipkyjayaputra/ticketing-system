package delivery

import (
	"github.com/sipkyjayaputra/ticketing-system/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Delivery interface {
	// AUTH
	SignIn(*gin.Context)
	RefreshToken(*gin.Context)

	// USERS
	GetUsers(*gin.Context)
	AddUser(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
	GetUserById(*gin.Context)

	// TICKET
	GetTickets(*gin.Context)
	AddTicket(*gin.Context)
	UpdateTicket(*gin.Context)
	DeleteTicket(*gin.Context)
	GetTicketById(*gin.Context)
}

type delivery struct {
	uc     usecase.Usecase
	logger *logrus.Logger
}

func NewDelivery(uc usecase.Usecase, logger *logrus.Logger) Delivery {
	return &delivery{
		uc:     uc,
		logger: logger,
	}
}
