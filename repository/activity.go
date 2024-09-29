package repository

import (
	"fmt"

	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"gorm.io/gorm"
)

// GetActivitiesByTicketNo retrieves all activities for a specific ticket from the database
func (repo *repository) GetActivitiesByTicketNo(ticketNo string) ([]entity.Activity, error) {
	activities := []entity.Activity{}
	if err := repo.db.Model(&entity.Activity{}).Preload("Documents").Where("ticket_no = ?", ticketNo).Order("created_at DESC").Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// AddActivity adds a new activity to the database
func (repo *repository) AddActivity(activity dto.Activity) error {
	// Begin the transaction
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Create the new activity entity
		newActivity := &entity.Activity{
			TicketNo:    activity.TicketNo, // Associate the activity with the provided ticketNo
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
			filePath := fmt.Sprintf("./uploads/%s/%s", activity.TicketNo, doc.Filename)

			// Save the uploaded file
			if err := helpers.SaveUploadedFile(doc, filePath); err != nil {
				return err
			}

			newDoc := entity.Document{
				ActivityID:   newActivity.ActivityID, // Use the generated ActivityID
				DocumentName: doc.Filename,
				DocumentSize: doc.Size,
				DocumentPath: filePath,
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
