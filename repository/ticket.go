package repository

import (
	"encoding/json"
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
		query = query.Where("DATE(report_date) <= ?", filter.ReportEndDate)
	}

	if filter.Terms != "" {
		query = query.Where("subject LIKE ?", "%"+filter.Terms+"%")
	}

	if filter.AssignedID != "" {
		query = query.Where("assigned_id = ?", filter.AssignedID)
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

	// Preload relationships, excluding `created_at` and `updated_at`
	userFields := "id, username" // Add other fields as needed
	if err := query.
		Preload("Activities.Documents").
		Preload("Updater", func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Preload("Assigned", func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Preload("Reporter", func(db *gorm.DB) *gorm.DB {
			return db.Select(userFields)
		}).
		Order("created_at DESC").
		Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (repo *repository) GetTicketSummary(ticket dto.TicketSummaryFilter) (*entity.TicketSummary, error) {
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
		newTicketCount, openTicketCount, inProgressTicketCount, closedTicketCount                                     int
		lastMonthNewTicketCount, lastMonthOpenTicketCount, lastMonthInProgressTicketCount, lastMonthClosedTicketCount int
	)

	// Function to build the query with or without assigned_id
	buildQuery := func(status string, month time.Month, year int) *gorm.DB {
		query := repo.db.Model(&entity.Ticket{}).
			Where("EXTRACT(MONTH FROM created_at) = ?", int(month)).
			Where("EXTRACT(YEAR FROM created_at) = ?", year)

		if ticket.Role != "admin" && ticket.Role != "management" {
			query = query.Where("assigned_id = ?", ticket.ID)
		}
		if status != "" {
			query = query.Where("UPPER(status) = ?", status)
		}

		return query
	}

	// Query for the current month
	if err := buildQuery("", currentMonth, currentYear).Select("COUNT(*) AS new_ticket_count").Scan(&newTicketCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery("OPEN", currentMonth, currentYear).Select("COUNT(*) AS open_ticket_count").Scan(&openTicketCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery("IN PROGRESS", currentMonth, currentYear).Select("COUNT(*) AS in_progress_ticket_count").Scan(&inProgressTicketCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery("CLOSED", currentMonth, currentYear).Select("COUNT(*) AS closed_ticket_count").Scan(&closedTicketCount).Error; err != nil {
		return nil, err
	}

	// Query for the last month
	if err := buildQuery("", lastMonth, lastYear).Select("COUNT(*) AS new_ticket_count").Scan(&lastMonthNewTicketCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery("OPEN", lastMonth, lastYear).Select("COUNT(*) AS open_ticket_count").Scan(&lastMonthOpenTicketCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery("IN PROGRESS", lastMonth, lastYear).Select("COUNT(*) AS in_progress_ticket_count").Scan(&lastMonthInProgressTicketCount).Error; err != nil {
		return nil, err
	}

	if err := buildQuery("CLOSED", lastMonth, lastYear).Select("COUNT(*) AS closed_ticket_count").Scan(&lastMonthClosedTicketCount).Error; err != nil {
		return nil, err
	}

	// Calculate notes
	ticketSummary.NewTicketCount = newTicketCount
	ticketSummary.NewTicketCountLastMonth = getTicketCountLastMonth(newTicketCount, lastMonthNewTicketCount)
	ticketSummary.NewTicketNote = getTicketNote(newTicketCount, lastMonthNewTicketCount)

	ticketSummary.OpenTicketCount = openTicketCount
	ticketSummary.OpenTicketCountLastMonth = getTicketCountLastMonth(openTicketCount, lastMonthOpenTicketCount)
	ticketSummary.OpenTicketNote = getTicketNote(openTicketCount, lastMonthOpenTicketCount)

	ticketSummary.InProgressTicketCount = inProgressTicketCount
	ticketSummary.InProgressTicketCountLastMonth = getTicketCountLastMonth(inProgressTicketCount, lastMonthInProgressTicketCount)
	ticketSummary.InProgressTicketNote = getTicketNote(inProgressTicketCount, lastMonthInProgressTicketCount)

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

func (repo *repository) AddTicket(ticket dto.Ticket) error {
    return repo.db.Transaction(func(tx *gorm.DB) error {
        // Get the last ticket ID
        var lastTicket entity.Ticket
        if err := tx.Order("ticket_id desc").First(&lastTicket).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
            return err // Rollback on error if not found
        }

        // Generate the new serial ID and ticket number
        newID := getNextTicketID(lastTicket.TicketNo)
        ticketNumber := utils.GenerateTicketNumber(uint(newID), ticket.TicketType)

        // Create the new ticket entity
        newTicket := entity.Ticket{
            TicketNo:   ticketNumber,
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
        }

        // Insert the new ticket
        if err := createEntity(tx, &newTicket); err != nil {
            return err
        }

        // Prepare to batch insert activities and documents
        activitiesToInsert := []entity.Activity{}
        documentsToInsert := []entity.Document{}

        for _, activity := range ticket.Activities {
            newActivity := entity.Activity{
                TicketID:    newTicket.TicketID,
                Description: "Initial Activity",
                CreatedBy:   ticket.CreatedBy,
                UpdatedBy:   ticket.UpdatedBy,
            }
            activitiesToInsert = append(activitiesToInsert, newActivity)

            // Prepare the documents for the activity
            for index, doc := range activity.Documents {
                filePath := fmt.Sprintf("./uploads/%s/%s", ticket.TicketType, doc.DocumentFile.Filename)

                if err := helpers.SaveUploadedFile(doc.DocumentFile, filePath); err != nil {
                    return err
                }

                var documentNo string

				var payload map[string]interface{}
                if err := json.Unmarshal(newTicket.Content, &payload); err == nil {
                    if originalPayload, ok := payload["original_payload"].(map[string]interface{}); ok {
                        if docNo, exists := originalPayload["document_no"].(string); exists {
                            documentNo = docNo
                        }
                    }
                }

                if documentNo == "" {
                    documentNo = utils.GenerateDocNumber(uint(newID), ticket.TicketType, index+1)
                }

                newDoc := entity.Document{
                    DocumentName: doc.DocumentFile.Filename,
                    DocumentNo:   documentNo,
                    DocumentSize: doc.DocumentFile.Size,
                    DocumentType: doc.DocumentType,
                    DocumentPath: filePath,
                    CreatedBy:    ticket.CreatedBy,
                    UpdatedBy:    ticket.UpdatedBy,
                }
                documentsToInsert = append(documentsToInsert, newDoc)
            }
        }

        // Bulk insert activities
        if len(activitiesToInsert) > 0 {
            if err := tx.Create(&activitiesToInsert).Error; err != nil {
                return err
            }
        }

        // Associate documents with activities
        for _, activity := range activitiesToInsert {
            for j := range documentsToInsert {
                documentsToInsert[j].ActivityID = activity.ActivityID // Update the ActivityID
            }
        }

        // Bulk insert documents
        if len(documentsToInsert) > 0 {
            if err := tx.Create(&documentsToInsert).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

func (repo *repository) UpdateTicket(updatedTicket dto.Ticket) error {
    return repo.db.Transaction(func(tx *gorm.DB) error {
        // Ambil tiket yang ingin diperbarui
        var existingTicket entity.Ticket
        if err := tx.First(&existingTicket, updatedTicket.TicketID).Error; err != nil {
            return err // Rollback jika tiket tidak ditemukan
        }

        // Update informasi tiket
        existingTicket.Subject = updatedTicket.Subject
        existingTicket.AssignedID = updatedTicket.AssignedID
        existingTicket.Priority = updatedTicket.Priority
        existingTicket.Status = updatedTicket.Status
        existingTicket.Content = updatedTicket.Content
        existingTicket.UpdatedBy = updatedTicket.UpdatedBy

        // Simpan perubahan tiket
        if err := tx.Save(&existingTicket).Error; err != nil {
            return err
        }

        // Ambil aktivitas pertama untuk mengupdate dokumen
        var activity entity.Activity
        if err := tx.Where("ticket_id = ?", updatedTicket.TicketID).First(&activity).Error; err != nil {
            return err // Rollback jika aktivitas tidak ditemukan
        }

        // Hapus dokumen yang ada
        if err := tx.Where("activity_id = ?", activity.ActivityID).Delete(&entity.Document{}).Error; err != nil {
            return err
        }

        // Menyimpan dokumen baru
        documentsToInsert := []entity.Document{}
        for index, doc := range updatedTicket.Activities[0].Documents {
            filePath := fmt.Sprintf("./uploads/%s/%s", updatedTicket.TicketType, doc.DocumentFile.Filename)

            if err := helpers.SaveUploadedFile(doc.DocumentFile, filePath); err != nil {
                return err
            }

            documentNo := utils.GenerateDocNumber(uint(updatedTicket.TicketID), updatedTicket.TicketType, index+1)

            newDoc := entity.Document{
                DocumentName: doc.DocumentFile.Filename,
                DocumentNo:   documentNo,
                DocumentSize: doc.DocumentFile.Size,
                DocumentType: doc.DocumentType,
                DocumentPath: filePath,
                CreatedBy:    updatedTicket.CreatedBy,
                UpdatedBy:    updatedTicket.UpdatedBy,
                ActivityID:   activity.ActivityID, // Kaitkan dokumen dengan aktivitas
            }
            documentsToInsert = append(documentsToInsert, newDoc)
        }

        // Bulk insert dokumen baru
        if len(documentsToInsert) > 0 {
            if err := tx.Create(&documentsToInsert).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

func (repo *repository) CloseTicket(ticket dto.Ticket) error {
	updateTicket := &entity.Ticket{
		Status: ticket.Status,
		UpdatedAt: ticket.UpdatedAt,
		UpdatedBy: ticket.UpdatedBy,
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
	userFields := "id, username"
	if err := repo.db.Model(&entity.Ticket{}).Preload("Activities.Updater").Preload("Activities.Documents").
	Preload("Updater", func(db *gorm.DB) *gorm.DB {
		return db.Select(userFields)
	}).
	Preload("Assigned", func(db *gorm.DB) *gorm.DB {
		return db.Select(userFields)
	}).
	Preload("Reporter", func(db *gorm.DB) *gorm.DB {
		return db.Select(userFields)
	}).Where("ticket_id = ?", id).First(&ticket).Error; err != nil {
		return nil, err
	}
	return ticket, nil
}

// Helper function to get the next ticket ID
func getNextTicketID(ticketNo string) int64 {
    lastDocNum, _ := strconv.ParseInt(strings.Split(ticketNo, "/")[0], 10, 64)
    return lastDocNum + 1
}

// Helper function to create an entity and handle errors
func createEntity(tx *gorm.DB, entity interface{}) error {
    return tx.Create(entity).Error
}