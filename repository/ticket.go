package repository

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"github.com/sipkyjayaputra/ticketing-system/utils"
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
	if err := query.Preload("Activities.Documents").Preload("Assigned").Preload("Reporter").Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// GetTicketSummary retrieves all tickets from the database
func (repo *repository) GetTicketSummary() (*entity.TicketSummary, error) {
	ticketSummary := entity.TicketSummary{}
	now := time.Now()
	currentMonth := now.Month()
	currentYear := now.Year()

	// Adjust last month and year
	lastMonth := currentMonth - 1
	lastYear := currentYear
	if lastMonth == 0 { // If the current month is January
		lastMonth = 12 // December
		lastYear--     // Move to last year
	}

	var (
		newTicketCount, openTicketCount, pendingTicketCount, closedTicketCount                                     int
		lastMonthNewTicketCount, lastMonthOpenTicketCount, lastMonthPendingTicketCount, lastMonthClosedTicketCount int
	)

	// Query for the current month
	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS new_ticket_count").
		Where("MONTH(created_at) = ?", currentMonth).
		Where("YEAR(created_at) = ?", currentYear).
		Scan(&newTicketCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS open_ticket_count").
		Where("UPPER(status) = ?", "OPEN").
		Where("MONTH(created_at) = ?", currentMonth).
		Where("YEAR(created_at) = ?", currentYear).
		Scan(&openTicketCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS pending_ticket_count").
		Where("UPPER(status) = ?", "PENDING").
		Where("MONTH(created_at) = ?", currentMonth).
		Where("YEAR(created_at) = ?", currentYear).
		Scan(&pendingTicketCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS closed_ticket_count").
		Where("UPPER(status) = ?", "CLOSED").
		Where("MONTH(created_at) = ?", currentMonth).
		Where("YEAR(created_at) = ?", currentYear).
		Scan(&closedTicketCount).Error; err != nil {
		return nil, err
	}

	// Query for the last month
	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS new_ticket_count").
		Where("MONTH(created_at) = ?", lastMonth).
		Where("YEAR(created_at) = ?", lastYear).
		Scan(&lastMonthNewTicketCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS open_ticket_count").
		Where("UPPER(status) = ?", "OPEN").
		Where("MONTH(created_at) = ?", lastMonth).
		Where("YEAR(created_at) = ?", lastYear).
		Scan(&lastMonthOpenTicketCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS pending_ticket_count").
		Where("UPPER(status) = ?", "PENDING").
		Where("MONTH(created_at) = ?", lastMonth).
		Where("YEAR(created_at) = ?", lastYear).
		Scan(&lastMonthPendingTicketCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&entity.Ticket{}).
		Select("COUNT(*) AS closed_ticket_count").
		Where("UPPER(status) = ?", "CLOSED").
		Where("MONTH(created_at) = ?", lastMonth).
		Where("YEAR(created_at) = ?", lastYear).
		Scan(&lastMonthClosedTicketCount).Error; err != nil {
		return nil, err
	}

	// Calculate notes
	ticketSummary.NewTicketCount = newTicketCount
	ticketSummary.NewTicketCountLastMonth = getTicketCountLastMonth(newTicketCount, lastMonthNewTicketCount)
	ticketSummary.NewTicketNote = getTicketNote(newTicketCount, lastMonthNewTicketCount)

	ticketSummary.OpenTicketCount = openTicketCount
	ticketSummary.OpenTicketCountLastMonth = getTicketCountLastMonth(openTicketCount, lastMonthOpenTicketCount)
	ticketSummary.OpenTicketNote = getTicketNote(openTicketCount, lastMonthOpenTicketCount)

	ticketSummary.PendingTicketCount = pendingTicketCount
	ticketSummary.PendingTicketCountLastMonth = getTicketCountLastMonth(pendingTicketCount, lastMonthPendingTicketCount)
	ticketSummary.PendingTicketNote = getTicketNote(pendingTicketCount, lastMonthPendingTicketCount)

	ticketSummary.ClosedTicketCount = closedTicketCount
	ticketSummary.ClosedTicketCountLastMonth = getTicketCountLastMonth(closedTicketCount, lastMonthClosedTicketCount)
	ticketSummary.ClosedTicketNote = getTicketNote(closedTicketCount, lastMonthClosedTicketCount)

	return &ticketSummary, nil
}

// Helper function to generate the note
func getTicketNote(currentCount, lastMonthCount int) string {
	if currentCount > lastMonthCount {
		return "+"
	} else if currentCount < lastMonthCount {
		return "-"
	} else {
		return ""
	}
}
func getTicketCountLastMonth(currentCount, lastMonthCount int) int {
	if currentCount > lastMonthCount {
		return currentCount -lastMonthCount
	} else if currentCount < lastMonthCount {
		return lastMonthCount - currentCount
	} else {
		return 0
	}
}

// AddTicket adds a new ticket to the database
func (repo *repository) AddTicket(ticket dto.Ticket) error {
	// Begin the transaction
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Get the last ticket ID
		var lastTicket entity.Ticket
		if err := tx.Order("ticket_id desc").First(&lastTicket).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err // Rollback on error if not found
			}
		}

		// Generate the new serial ID (last ID + 1)
		newID := lastTicket.TicketID + 1

		// Generate the ticket number
		ticketNumber := utils.GenerateTicketNumber(newID, ticket.TicketType)

		// Create the new ticket entity
		newTicket := &entity.Ticket{
			TicketNo:   ticketNumber, // Use the generated ticket number
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
				TicketID:    newID, // Associate the activity with the created ticket
				Description: "Initial Activity",
				CreatedBy:   ticket.CreatedBy,
				UpdatedBy:   ticket.UpdatedBy,
				CreatedAt:   ticket.CreatedAt,
				UpdatedAt:   ticket.UpdatedAt,
			}

			// Insert the activity
			if err := tx.Create(&newActivity).Error; err != nil {
				return err // Rollback on error
			}

			// Prepare the documents for the activity
			for index, doc := range activity.Documents {
				filePath := fmt.Sprintf("./uploads/%s/%s", ticket.TicketType, doc.DocumentFile.Filename)

				if err := helpers.SaveUploadedFile(doc.DocumentFile, filePath); err != nil {
					return err
				}

				newDoc := entity.Document{
					ActivityID:   newActivity.ActivityID,
					DocumentName: doc.DocumentFile.Filename,
					DocumentNo: utils.GenerateDocNumber(newID,  ticket.TicketType, index+1),
					DocumentSize: doc.DocumentFile.Size,
					DocumentType: doc.DocumentType,
					DocumentPath: filePath,
					CreatedBy:    ticket.CreatedBy,
					UpdatedBy:    ticket.UpdatedBy,
					CreatedAt:    ticket.CreatedAt,
					UpdatedAt:    ticket.UpdatedAt,
				}
				// Insert the document
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

	if err := repo.db.Model(&entity.Ticket{}).Where("ticket_id = ?", ticket.TicketID).Updates(&updateTicket).Error; err != nil {
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
	if err := repo.db.Model(&entity.Ticket{}).Preload("Activities").Preload("Activities.Documents").Preload("Assigned").Preload("Reporter").Where("ticket_id = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}
