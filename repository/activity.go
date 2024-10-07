package repository

import (
	"errors"
	"fmt"

	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"gorm.io/gorm"
)

// GetActivitiesByTicketNo retrieves all activities for a specific ticket from the database
func (repo *repository) GetActivitiesByTicketNo(id string) ([]entity.Activity, error) {
	activities := []entity.Activity{}
	if err := repo.db.Model(&entity.Activity{}).Preload("Documents").Where("id = ?", id).Order("created_at DESC").Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// AddActivity adds a new activity to the database
func (repo *repository) AddActivity(activity dto.Activity) error {

	// Begin the transaction
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// UPDATE TICKET STATUS
		activities := []entity.Activity{}

		if err := repo.db.Model(&entity.Activity{}).Where("ticket_id = ?", activity.TicketID).Find(&activities).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err // Rollback on error if not found
			}
		}

		if len(activities) == 1 {
			// Update ticket status into In Progress
			if err := tx.Model(&entity.Ticket{}).Where("ticket_id = ?", activity.TicketID).Update("status", "In Progress").Error; err != nil {
				return err
			}
		}
		// UPDATE TICKET STATUS

		// Create the new activity entity
		newActivity := &entity.Activity{
			TicketID:    activity.TicketID, // Associate the activity with the provided ticketNo
			Description: activity.Description,
			CreatedBy:   activity.CreatedBy,
			UpdatedBy:   activity.UpdatedBy,
			CreatedAt:   activity.CreatedAt,
			UpdatedAt:   activity.UpdatedAt,
		}

		// Insert the new activity
		if err := tx.Create(&newActivity).Error; err != nil {
			return err // Rollback on error
		}

		// Loop through each document related to the activity
		for _, doc := range activity.Documents {
			filePath := fmt.Sprintf("./uploads/%d/%s", activity.TicketID, doc.DocumentFile.Filename)

			// Save the uploaded file
			if err := helpers.SaveUploadedFile(doc.DocumentFile, filePath); err != nil {
				return err
			}

			newDoc := entity.Document{
				ActivityID:   newActivity.ActivityID, // Use the generated ActivityID
				DocumentName: doc.DocumentFile.Filename,
				DocumentSize: doc.DocumentFile.Size,
				DocumentPath: filePath,
				DocumentType: doc.DocumentType,
				CreatedBy:    activity.CreatedBy,
				UpdatedBy:    activity.UpdatedBy,
				CreatedAt:    activity.CreatedAt,
				UpdatedAt:    activity.UpdatedAt,
			}

			// Insert the document
			if err := tx.Create(&newDoc).Error; err != nil {
				return err // Rollback on error
			}
		}

		// Commit the transaction
		return nil
	})
}

// UpdateActivity updates an existing activity in the database
func (repo *repository) UpdateActivity(activity dto.Activity) error {
	updateActivity := &entity.Activity{
		Description: activity.Description,
		UpdatedBy:   activity.UpdatedBy,
		UpdatedAt:   activity.UpdatedAt,
	}

	if err := repo.db.Model(&entity.Activity{}).Where("activity_id = ?", activity.ActivityID).Updates(&updateActivity).Error; err != nil {
		return err
	}
	return nil
}

// DeleteActivity deletes an activity from the database
func (repo *repository) DeleteActivity(activityID uint) error {
	if err := repo.db.Where("activity_id = ?", activityID).Delete(&entity.Activity{}).Error; err != nil {
		return err
	}
	return nil
}

// GetActivityById retrieves a specific activity by ID from the database
func (repo *repository) GetActivityById(activityID uint) (*entity.Activity, error) {
	activity := &entity.Activity{}
	if err := repo.db.Model(&entity.Activity{}).Preload("Documents").Where("activity_id = ?", activityID).First(&activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}
