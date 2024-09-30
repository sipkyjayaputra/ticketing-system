package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"gorm.io/gorm"
)

// GetTickets retrieves all tickets from the database
func (repo *repository) GetTickets(filter dto.TicketFilter) ([]entity.Ticket, error) {
	tickets := []entity.Ticket{}
	query := repo.db.Model(&entity.Ticket{})

	// Apply filters
	if filter.TicketType != "" {
		query = query.Where("UPPER(ticket_type) = ?", strings.ToUpper(filter.TicketType))
	}

	if filter.Priority != "" {
		query = query.Where("UPPER(priority) = ?", strings.ToUpper(filter.Priority))
	}

	if filter.Status != "" {
		query = query.Where("UPPER(status) = ?", strings.ToUpper(filter.Status))
	}

	if filter.ReportStartDate != "" && filter.ReportEndDate != "" {
		query = query.Where("DATE(report_date) BETWEEN ? AND ?", filter.ReportStartDate, filter.ReportEndDate)
	} else if filter.ReportStartDate != "" {
		query = query.Where("DATE(report_date) >= ?", filter.ReportStartDate)
	} else if filter.ReportEndDate != "" {
		query = query.Where("DATE(report_date) <= ?", filter.ReportStartDate)
	}

	if filter.Terms != "" {
		query = query.Where("subject LIKE ?", "%"+filter.Terms+"%")
	}

	// Handle pagination (limit and offset)
	if filter.Limit != "" {
		limit, err := strconv.ParseInt(filter.Limit, 10, 64)
		if err == nil && limit > 0 {
			query = query.Limit(int(limit))
		}
	}

	if filter.Offset != "" {
		offset, err := strconv.ParseInt(filter.Offset, 10, 64)
		if err == nil && offset >= 0 {
			query = query.Offset(int(offset))
		}
	}

	// Preload relationships and order by creation date
	if err := query.Preload("Activities").Preload("Assigned").Preload("Reporter").Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// AddTicket adds a new ticket to the database
func (repo *repository) AddTicket(ticket dto.Ticket) error {
	// Begin the transaction
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Create the new ticket entity
		newTicket := &entity.Ticket{
			TicketNo:   ticket.TicketNo,
			TicketType: ticket.TicketType,
			Subject:    ticket.Subject,
			ReportDate: ticket.ReportDate,
			AssignedID: ticket.AssignedID,
			ReporterID: ticket.ReporterID,
			Priority:   ticket.Priority,
			Status:     ticket.Status,
			Content:    ticket.Content,
			CreatedBy:  ticket.CreatedBy,
			UpdatedBy:  ticket.UpdatedBy,
			CreatedAt:  ticket.CreatedAt,
			UpdatedAt:  ticket.UpdatedAt,
		}

		// Insert the new ticket
		if err := tx.Create(&newTicket).Error; err != nil {
			return err // Rollback on error
		}

		// Loop through each activity related to the ticket
		for _, activity := range ticket.Activities {
			// Create the new activity entity
			newActivity := entity.Activity{
				TicketNo:    newTicket.TicketNo, // Associate the activity with the created ticket
				Description: "Initial Activity",
				CreatedBy:   ticket.CreatedBy,
				UpdatedBy:   ticket.UpdatedBy,
				CreatedAt:   ticket.CreatedAt,
				UpdatedAt:   ticket.UpdatedAt,
			}

			// Insert the activity along with its documents
			if err := tx.Create(&newActivity).Error; err != nil {
				return err // Rollback on error
			}

			// Prepare the documents for the activity
			for _, doc := range activity.Documents {
				filePath := fmt.Sprintf("./uploads/%s/%s", ticket.TicketType, doc.Filename)

				if err := helpers.SaveUploadedFile(doc, filePath); err != nil {
					return err
				}

				newDoc := entity.Document{
					ActivityID:   newActivity.ActivityID,
					DocumentName: doc.Filename,
					DocumentSize: doc.Size,
					DocumentPath: filePath,
					CreatedBy:    ticket.CreatedBy,
					UpdatedBy:    ticket.UpdatedBy,
					CreatedAt:    ticket.CreatedAt,
					UpdatedAt:    ticket.UpdatedAt,
				}
				// Insert the activity along with its documents
				if err := tx.Create(&newDoc).Error; err != nil {
					return err // Rollback on error
				}
			}
		}

		// Commit the transaction
		return nil
	})
}

// UpdateTicket updates an existing ticket in the database
func (repo *repository) UpdateTicket(ticket dto.Ticket) error {
	updateTicket := &entity.Ticket{
		TicketNo:   ticket.TicketNo,
		TicketType: ticket.TicketType,
		Subject:    ticket.Subject,
		ReportDate: ticket.ReportDate,
		AssignedID: ticket.AssignedID,
		Priority:   ticket.Priority,
		Status:     ticket.Status,
		Content:    ticket.Content,
		CreatedBy:  ticket.CreatedBy,
		UpdatedBy:  ticket.UpdatedBy,
		CreatedAt:  ticket.CreatedAt,
		UpdatedAt:  ticket.UpdatedAt,
	}

	if err := repo.db.Model(&entity.Ticket{}).Where("ticket_no = ?", ticket.TicketNo).Updates(&updateTicket).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTicket deletes a ticket from the database
func (repo *repository) DeleteTicket(id string) error {
	if err := repo.db.Where("ticket_no = ?", id).Delete(&entity.Ticket{}).Error; err != nil {
		return err
	}
	return nil
}

// GetTicketById retrieves a specific ticket by ID from the database
func (repo *repository) GetTicketById(id string) (*entity.Ticket, error) {
	ticket := &entity.Ticket{}
	if err := repo.db.Model(&entity.Ticket{}).Preload("Activities").Preload("Activities.Documents").Preload("Assigned").Preload("Reporter").Where("ticket_no = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}
