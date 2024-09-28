package delivery

import (
	"github.com/sipkyjayaputra/ticketing-system/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Delivery interface {
	SignIn(*gin.Context)
	RefreshToken(*gin.Context)
	GetUsers(*gin.Context)
	AddUser(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
	GetUserById(*gin.Context)
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
