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
	Role      string     `gorm:"column:role" json:"role"`
	Photo     string     `gorm:"column:photo;default:''" json:"photo"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	CreatedBy string     `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedBy string     `gorm:"column:updated_by" json:"updated_by,omitempty"`
}

type Document struct {
	DocumentID   uint      `gorm:"primaryKey;autoIncrement" json:"document_id"`
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
	TicketNo    string    `gorm:"column:ticket_no" json:"ticket_no"`     // Foreign key, references Ticket's TicketNo
	Description string    `gorm:"column:description" json:"description"` // Description of the activity
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`      // Activity creation timestamp
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // Activity update timestamp
	CreatedBy   uint      `gorm:"column:created_by" json:"created_by"`   // User who created the activity
	UpdatedBy   uint      `gorm:"column:updated_by" json:"updated_by"`   // User who last updated the activity

	// Relationships
	Documents []Document `gorm:"foreignKey:ActivityID" json:"documents"` // Multiple documents related to the activity
}

type Ticket struct {
	TicketNo   string          `gorm:"primaryKey;column:ticket_no" json:"ticket_no"` // Unique ticket number
	ReporterID uint            `gorm:"column:reporter_id" json:"reporter_id"`        // ID of the user who reported the ticket
	TicketType string          `gorm:"column:ticket_type" json:"ticket_type"`        // Type of ticket
	Subject    string          `gorm:"column:subject" json:"subject"`                // Subject of the ticket
	ReportDate time.Time       `gorm:"column:report_date" json:"report_date"`        // Date the ticket was reported
	AssignedID uint            `gorm:"column:assigned_id" json:"assigned_id"`        // ID of the user assigned to the ticket
	Priority   string          `gorm:"column:priority" json:"priority"`              // Priority of the ticket
	Status     string          `gorm:"column:status" json:"status"`                  // Status of the ticket
	Content    json.RawMessage `gorm:"column:content" json:"content"`                // Raw JSON content of the ticket
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`             // Ticket creation timestamp
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`             // Ticket update timestamp
	CreatedBy  uint            `gorm:"column:created_by" json:"created_by"`          // User who created the ticket
	UpdatedBy  uint            `gorm:"column:updated_by" json:"updated_by"`          // User who last updated the ticket

	// Relationships
	Assigned   User       `gorm:"foreignKey:AssignedID" json:"assigned"` // Assigned user relationship
	Reporter   User       `gorm:"foreignKey:ReporterID" json:"reporter"` // Reporter user relationship
	Activities []Activity `gorm:"foreignKey:TicketNo" json:"activities"` // Multiple activities related to the ticket
}

type TicketSummary struct {
	NewTicketCount     int    `json:"new_ticket_count"`
	NewTicketNote      string `json:"new_ticket_note"`
	OpenTicketCount    int    `json:"open_ticket_count"`
	OpenTicketNote     string `json:"open_ticket_note"`
	PendingTicketCount int    `json:"pending_ticket_count"`
	PendingTicketNote  string `json:"pending_ticket_note"`
	ClosedTicketCount  int    `json:"closed_ticket_count"`
	ClosedTicketNote   string `json:"closed_ticket_note"`
}
