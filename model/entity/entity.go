package entity

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        uint       `gorm:"primaryKey;autoIncrement;index" json:"id"`
	Username  string     `gorm:"column:username;unique" json:"username"`    // Username harus unik
	Password  string     `gorm:"column:password" json:"password,omitempty"` // Password pengguna
	Email     string     `gorm:"column:email;unique" json:"email"`          // Email pengguna, harus unik
	Phone     string     `gorm:"column:phone;unique" json:"phone"`          // Email pengguna, harus unik
	Team     string     `gorm:"column:team" json:"team"`          // Email pengguna, harus unik
	Workpalce     string     `gorm:"column:workplace" json:"workplace"`          // Email pengguna, harus unik
	Role      string     `gorm:"column:role" json:"role"`
	Photo     string     `gorm:"column:photo;default:''" json:"photo"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	CreatedBy string     `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy string     `gorm:"column:updated_by" json:"updated_by,omitempty"`
}

type Document struct {
	DocumentID   uint      `gorm:"primaryKey;autoIncrement" json:"document_id"`
	DocumentNo   string      `gorm:"column:document_no" json:"document_no"`
	ActivityID   uint      `gorm:"column:activity_id" json:"activity_id"`     // Foreign key, references Activity's ActivityID
	DocumentName string    `gorm:"column:document_name" json:"document_name"` // Name of the document
	DocumentSize int64     `gorm:"column:document_size" json:"document_size"` // Size of the document in bytes
	DocumentPath string    `gorm:"column:document_path" json:"document_path"` // Path for the document
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`    // Document creation timestamp
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`    // Document update timestamp
	CreatedBy    uint      `gorm:"column:created_by" json:"created_by"` // User who uploaded the document
	UpdatedBy    uint      `gorm:"column:updated_by" json:"updated_by"` // User who last updated the document
}

type Activity struct {
	ActivityID  uint      `gorm:"primaryKey;autoIncrement" json:"activity_id"`
	TicketID    uint      `gorm:"column:ticket_id" json:"ticket_id"`     // Foreign key, references Ticket's TicketNo
	Description string    `gorm:"column:description" json:"description"` // Description of the activity
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`      // Activity creation timestamp
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // Activity update timestamp
	CreatedBy   uint      `gorm:"column:created_by" json:"created_by"`   // User who created the activity
	UpdatedBy   uint      `gorm:"column:updated_by" json:"updated_by"`   // User who last updated the activity

	// Relationships
	Documents []Document `gorm:"foreignKey:ActivityID" json:"documents"` // Multiple documents related to the activity
	Updater    User      `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`  
}


type Ticket struct {
	TicketID   uint            `gorm:"primaryKey;autoIncrement" json:"ticket_id"`      // Serial ID for the ticket
	TicketNo   string          `gorm:"column:ticket_no" json:"ticket_no"`               // Unique ticket number
	ReporterID uint            `gorm:"column:reporter_id" json:"reporter_id"`           // ID of the user who reported the ticket
	TicketType string          `gorm:"column:ticket_type" json:"ticket_type"`           // Type of ticket
	Subject    string          `gorm:"column:subject" json:"subject"`                   // Subject of the ticket
	ReportDate time.Time       `gorm:"column:report_date" json:"report_date"`           // Date the ticket was reported
	AssignedID uint            `gorm:"column:assigned_id" json:"assigned_id"`           // ID of the user assigned to the ticket
	Priority   string          `gorm:"column:priority" json:"priority"`                 // Priority of the ticket
	Status     string          `gorm:"column:status" json:"status"`                     // Status of the ticket
	Content    json.RawMessage `gorm:"column:content" json:"content"`                   // Raw JSON content of the ticket
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`                // Ticket creation timestamp
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`                // Ticket update timestamp
	CreatedBy  uint            `gorm:"column:created_by" json:"created_by"`             // User who created the ticket
	UpdatedBy  uint            `gorm:"column:updated_by" json:"updated_by"`             // User who last updated the ticket

	// Relationships
	Assigned   User       `gorm:"foreignKey:AssignedID;references:ID" json:"assigned"`   // Assigned user relationship
	Reporter   User       `gorm:"foreignKey:ReporterID;references:ID" json:"reporter"`   // Reporter user relationship
	Updater    User       `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`     // User who last updated the ticket
	Activities []Activity  `gorm:"foreignKey:TicketID" json:"activities"`    // Multiple activities related to the ticket
}

type TicketSummary struct {
	NewTicketCount     int    `json:"new_ticket_count"`
	NewTicketCountLastMonth int `json:"new_ticket_count_last_month"`
	NewTicketNote string `json:"new_ticket_note"`
	OpenTicketCount    int    `json:"open_ticket_count"`
	OpenTicketCountLastMonth     int `json:"open_ticket_count_last_month"`
	OpenTicketNote     string `json:"open_ticket_note"`
	PendingTicketCount int    `json:"pending_ticket_count"`
	PendingTicketCountLastMonth  int `json:"pending_ticke_count_last_month"`
	PendingTicketNote  string `json:"pending_ticket_note"`
	ClosedTicketCount  int    `json:"closed_ticket_count"`
	ClosedTicketCountLastMonth   int `json:"closed_ticket_count_last_month"`
	ClosedTicketNote   string `json:"closed_ticket_note"`
}
