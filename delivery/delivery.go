package delivery

import (
	"github.com/sipkyjayaputra/ticketing-system/helpers"
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
	UpdateUserPhoto(*gin.Context)
	UpdateUserPassword(*gin.Context)

	// TICKET
	GetTickets(*gin.Context)
	AddTicket(*gin.Context)
	UpdateTicket(*gin.Context)
	DeleteTicket(*gin.Context)
	GetTicketById(*gin.Context)
	GetTicketSummary(*gin.Context)
	CloseTicket(*gin.Context)

	// ACTIVITY
	GetActivitiesByTicketNo(*gin.Context)
	AddActivity(*gin.Context)
	UpdateActivity(*gin.Context)
	DeleteActivity(*gin.Context)
	GetActivityById(*gin.Context)

	// FILE
	FileServe(*gin.Context)
	FileServeByPath(*gin.Context)
	FileDownload(*gin.Context)

	// HRSV
	GetHrsvUsers(*gin.Context)
	GetHrsvRoles(*gin.Context)
	SyncUserDataHrsv(*gin.Context)
	SyncPasswordHrsv(*gin.Context)
}

type delivery struct {
	uc     usecase.Usecase
	logger *logrus.Logger
	hrsvClient *helpers.HrsvClient
}

func NewDelivery(uc usecase.Usecase, logger *logrus.Logger, hrsvClient *helpers.HrsvClient) Delivery {
	return &delivery{
		uc:     uc,
		logger: logger,
		hrsvClient: hrsvClient,
	}
}
