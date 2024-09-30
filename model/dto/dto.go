package dto

import (
	"encoding/json"
	"mime/multipart"
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

type TicketFilter struct {
	TicketType      string `json:"ticket_type"`
	Priority        string `json:"priority"`
	Status          string `json:"status"`
	ReportStartDate string `json:"report_start_date"`
	ReportEndDate   string `json:"report_end_date"`
	Terms           string `json:"terms"`
	Limit           string `json:"limit"`
	Offset          string `json:"offset"`
}

type Ticket struct {
	TicketNo   string          `json:"ticket_no"`
	ReporterID uint            `json:"reporter_id"`
	TicketType string          `json:"ticket_type"`
	Subject    string          `json:"subject"`
	ReportDate time.Time       `json:"report_date"`
	AssignedID uint            `json:"assigned_id"`
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
	ActivityID  uint                    `json:"activity_id"`
	TicketNo    string                  `json:"ticket_no"`
	Description string                  `json:"description"`
	Documents   []*multipart.FileHeader `json:"documents"`
	CreatedBy   uint                    `json:"created_by,omitempty"` // ID of the creator
	UpdatedBy   uint                    `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt   time.Time               `json:"created_at"`           // Activity creation timestamp
	UpdatedAt   time.Time               `json:"updated_at"`           // Activity update timestamp
}

type Document struct {
	DocumentID   uint      `json:"document_id"`
	DocumentName string    `json:"document_name"`
	DocumentSize int64     `json:"document_size"`
	DocumentPath string    `json:"document_path"`
	CreatedBy    uint      `json:"created_by,omitempty"` // ID of the uploader
	UpdatedBy    uint      `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt    time.Time `json:"created_at"`           // Document creation timestamp
	UpdatedAt    time.Time `json:"updated_at"`           // Document update timestamp
}
