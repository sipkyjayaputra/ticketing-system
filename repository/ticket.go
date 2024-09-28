package repository

import (
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"gorm.io/gorm"
)

// GetTickets retrieves all tickets from the database
func (repo *repository) GetTickets() ([]entity.Ticket, error) {
	tickets := []entity.Ticket{}
	if err := repo.db.Model(&entity.Ticket{}).Order("created_at DESC").Find(&tickets).Error; err != nil {
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
			AssignedID: ticket.Assigned.ID,
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
				Description: activity.Description,
				CreatedAt:   activity.CreatedAt,
				UpdatedAt:   activity.UpdatedAt,
				CreatedBy:   activity.CreatedBy,
				UpdatedBy:   activity.UpdatedBy,
			}

			// Insert the activity along with its documents
			if err := tx.Create(&newActivity).Error; err != nil {
				return err // Rollback on error
			}

			// Prepare the documents for the activity
			for _, doc := range activity.Files {
				newDoc := entity.Document{
					ActivityID:   activity.ActivityID,
					DocumentName: doc.DocumentName,
					DocumentSize: doc.DocumentSize,
					DocumentBlob: doc.DocumentBlob,
					CreatedAt:    doc.CreatedAt,
					UpdatedAt:    doc.UpdatedAt,
					CreatedBy:    doc.CreatedBy,
					UpdatedBy:    doc.UpdatedBy,
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
func (repo *repository) UpdateTicket(ticket dto.Ticket, id string) error {
	updateTicket := &entity.Ticket{
		TicketNo:   ticket.TicketNo,
		TicketType: ticket.TicketType,
		Subject:    ticket.Subject,
		ReportDate: ticket.ReportDate,
		AssignedID: ticket.Assigned.ID,
		Priority:   ticket.Priority,
		Status:     ticket.Status,
		Content:    ticket.Content,
		CreatedBy:  ticket.CreatedBy,
		UpdatedBy:  ticket.UpdatedBy,
		CreatedAt:  ticket.CreatedAt,
		UpdatedAt:  ticket.UpdatedAt,
	}

	if err := repo.db.Model(&entity.Ticket{}).Where("ticket_id = ?", id).Updates(&updateTicket).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTicket deletes a ticket from the database
func (repo *repository) DeleteTicket(id string) error {
	if err := repo.db.Where("ticket_id = ?", id).Delete(&entity.Ticket{}).Error; err != nil {
		return err
	}
	return nil
}

// GetTicketById retrieves a specific ticket by ID from the database
func (repo *repository) GetTicketById(id string) (*entity.Ticket, error) {
	ticket := &entity.Ticket{}
	if err := repo.db.Model(&entity.Ticket{}).Where("ticket_id = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}
