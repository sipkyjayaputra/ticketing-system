package dto

import (
	"encoding/json"
	"mime/multipart"
	"time"
)

type UpdateUserPhoto struct {
	ID    uint                  `json:"id"`
	Photo *multipart.FileHeader `json:"photo"`
}
type CloseTicket struct {
	TicketID    uint   `json:"ticket_id"`
	Status 		string `json:"status"`
}

type UpdateUserPassword struct {
	ID             string   `json:"id"`
	CurrentPassword string `json:"current_password"`
	NewPassword    string `json:"new_password"`
	VerifyPassword string `json:"verify_password"`
}

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username" binding:"required"` // Username harus unik
	Password  string    `json:"password,omitempty"`          // Password pengguna
	Email     string    `json:"email" binding:"required"`    // Email pengguna, harus unik
	Role      string    `json:"role" binding:"required"`     // ID peran yang terkait
	Phone     string    `json:"phone"`     // ID peran yang terkait
	Workpalce string    `json:"workplace" `     // ID peran yang terkait
	Team      string    `json:"team"`     // ID peran yang terkait
	CreatedBy string    `json:"created_by,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SignIn struct {
	Email string `json:"email" binding:"required"` // Username harus unik
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
	AssignedID 		string `json:"assigned_id"`
}

type Ticket struct {
	TicketID   uint            `json:"ticket_id"`
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
	ActivityID  uint       `json:"activity_id"`
	TicketID    uint       `json:"ticket_id"` 			// Ticket
	Description string     `json:"description"`
	Documents   []Document `json:"documents"`
	CreatedBy   uint       `json:"created_by,omitempty"` // ID of the creator
	UpdatedBy   uint       `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt   time.Time  `json:"created_at"`           // Activity creation timestamp
	UpdatedAt   time.Time  `json:"updated_at"`           // Activity update timestamp
}

type Document struct {
	DocumentID   uint                  `json:"document_id"`
	DocumentNo   string                `json:"document_no"`
	DocumentName string                `json:"document_name"`
	DocumentSize int64                 `json:"document_size"`
	DocumentPath string                `json:"document_path"`
	DocumentType string                `json:"document_type"`
	DocumentFile *multipart.FileHeader `json:"document_file"`
	CreatedBy    uint                  `json:"created_by,omitempty"` // ID of the uploader
	UpdatedBy    uint                  `json:"updated_by,omitempty"` // ID of the last updater
	CreatedAt    time.Time             `json:"created_at"`           // Document creation timestamp
	UpdatedAt    time.Time             `json:"updated_at"`           // Document update timestamp
}


type Response struct {
    StatusCode      int         `json:"status_code"`
    Success         bool        `json:"success"`
    ResponseMessage string      `json:"response_message"`
    Errors          interface{} `json:"errors"` // Use interface{} if errors can be null or an object
    Data            []map[string]interface{}  `json:"data"`
}

// UserData struct to hold user details
type UserDataHRSV struct {
    ID           string       `json:"_id"`
    CompanyRole  CompanyRole  `json:"company_role"`
    CreatedAt    time.Time    `json:"created_at"`
    Email        string       `json:"email"`
    FID          string       `json:"fid"`
    IsActive     bool         `json:"is_active"`
    IsSynced     bool         `json:"is_synced"`
    IsVerified    bool        `json:"is_verified"`
    Name         string       `json:"name"`
    ProfileDesc  string       `json:"profile_desc"`
    SystemRole   string       `json:"system_role"`
    Team         Team         `json:"team"`
    UpdatedAt    time.Time    `json:"updated_at"`
    Workplace    Workplace    `json:"workplace"`
}

// CompanyRole struct to hold company role details
type CompanyRole struct {
    FID  string `json:"fid"`
    Name string `json:"name"`
}

// Team struct to hold team details
type Team struct {
    FID  string `json:"fid"`
    Name string `json:"name"`
}

// Workplace struct to hold workplace details
type Workplace struct {
    FID     string    `json:"fid"`
    Location Location  `json:"location"`
    Name    string     `json:"name"`
    Radius  int       `json:"radius"`
}

// Location struct to hold location details
type Location struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

type TicketSummaryFilter struct {
	ID uint `json:"id"`
	Role string `json:"role"`
}