package usecase

import (
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"
)

// GetActivitiesByTicketNo retrieves all activities related to a specific ticket
func (uc *usecase) GetActivitiesByTicketNo(ticketNo string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	activities, err := uc.repo.GetActivitiesByTicketNo(ticketNo)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get activities", err.Error())
	}

	return utils.BuildSuccessResponse(activities), nil
}

// AddActivity adds a new activity for a specific ticket
func (uc *usecase) AddActivity(activity dto.Activity, creator uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	activity.CreatedBy = creator
	activity.UpdatedBy = creator
	activity.CreatedAt = time.Now()
	activity.UpdatedAt = time.Now()

	if err := uc.repo.AddActivity(activity); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to add activity", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

// UpdateActivity updates an existing activity
func (uc *usecase) UpdateActivity(activity dto.Activity, updater uint, activityID uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	activity.ActivityID = activityID
	activity.UpdatedBy = updater
	activity.UpdatedAt = time.Now()

	if err := uc.repo.UpdateActivity(activity); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to update activity", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

// DeleteActivity deletes an activity by its ID
func (uc *usecase) DeleteActivity(activityID uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	if err := uc.repo.DeleteActivity(activityID); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to delete activity", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

// GetActivityById retrieves a specific activity by its ID
func (uc *usecase) GetActivityById(activityID uint) (*utils.ResponseContainer, *utils.ErrorContainer) {
	activity, err := uc.repo.GetActivityById(activityID)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get activity", err.Error())
	}

	return utils.BuildSuccessResponse(activity), nil
}
