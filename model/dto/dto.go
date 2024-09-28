package dto

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username" binding:"required"` // Username harus unik
	Password  string    `json:"password,omitempty"`          // Password pengguna
	Email     string    `json:"email" binding:"required"`    // Email pengguna, harus unik
	Role      string    `json:"role" binding:"required"`     // ID peran yang terkait
	CreatedBy string    `json:"created_by,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SignIn struct {
	Username string `json:"username" binding:"required"` // Username harus unik
	Password string `json:"password" binding:"required"` // Password pengguna
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Ticket struct {
	TicketNo   string          `json:"ticket_no"`
	Reporter   User            `json:"reporter"`
	TicketType string          `json:"ticket_type"`
	Subject    string          `json:"subject"`
	ReportDate time.Time       `json:"report_date"`
	Assigned   User            `json:"assigned"`
	Priority   string          `json:"priority"`
	Status     string          `json:"status"`
	Content    json.RawMessage `json:"content"`
	CreatedBy  uint            `json:"created_by,omitempty"` // ID of the creator
	UpdatedBy  uint            `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt  time.Time       `json:"created_at"`           // Ticket creation timestamp
	UpdatedAt  time.Time       `json:"updated_at"`           // Ticket update timestamp
	Activities []Activity      `json:"activities"`           // Associated activities
}

type Activity struct {
	ActivityID  uint       `json:"activity_id"`
	TicketNo    string     `json:"ticket_no"`
	Description string     `json:"description"`
	Files       []Document `json:"files"`
	CreatedBy   uint       `json:"created_by,omitempty"` // ID of the creator
	UpdatedBy   uint       `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt   time.Time  `json:"created_at"`           // Activity creation timestamp
	UpdatedAt   time.Time  `json:"updated_at"`           // Activity update timestamp
}

type Document struct {
	DocumentID   uint      `json:"document_id"`
	DocumentName string    `json:"document_name"`
	DocumentSize int64     `json:"document_size"`
	DocumentBlob string    `json:"document_blob"`
	CreatedBy    uint      `json:"created_by,omitempty"` // ID of the uploader
	UpdatedBy    uint      `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt    time.Time `json:"created_at"`           // Document creation timestamp
	UpdatedAt    time.Time `json:"updated_at"`           // Document update timestamp
}
