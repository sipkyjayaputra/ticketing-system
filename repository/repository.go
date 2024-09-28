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

	// TICKET
	// GetTickets() ([]entity.Ticket, error)
	AddTicket(dto.Ticket) error
	// UpdateTicket(dto.Ticket) error
	// DeleteTicket(uint) error
	// GetTicketById(uint) (*entity.Ticket, error)
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