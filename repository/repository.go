package repository

import (
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository interface {
	// AUTH
	SignIn(string) (entity.User, error)

	// USER
	GetUsers() ([]entity.User, error)
	AddUser(dto.User) error
	UpdateUser(dto.User, string) error
	DeleteUser(string) error
	GetUserById(string) (*entity.User, error)
	GetUserPassword(string) (string, error)
	UpdateUserPassword(dto.UpdateUserPassword) error
	UpdateUserPhoto(dto.UpdateUserPhoto) error

	// TICKET
	GetTickets(dto.TicketFilter) ([]entity.Ticket, error)
	AddTicket(dto.Ticket) error
	UpdateTicket(dto.Ticket) error
	CloseTicket(dto.Ticket) error
	DeleteTicket(string) error
	GetTicketById(string) (*entity.Ticket, error)
	GetTicketSummary(dto.TicketSummaryFilter) (*entity.TicketSummary, error)

	// ACTIVITY
	GetActivitiesByTicketNo(string) ([]entity.Activity, error)
	AddActivity(dto.Activity) error
	UpdateActivity(dto.Activity) error
	DeleteActivity(uint) error
	GetActivityById(uint) (*entity.Activity, error)

	// DOCUMENT
	GetDocumentById(uint) (*entity.Document, error)

	// HRSV
	SyncUserDataHrsv([]dto.UserDataHRSV) error
}

type repository struct {
	db     *gorm.DB
	cache  *redis.Client
	logger *logrus.Logger
}

func NewRepository(db *gorm.DB, logger *logrus.Logger, cache *redis.Client) Repository {
	return &repository{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}
