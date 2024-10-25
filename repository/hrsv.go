package repository

import (
	"fmt"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
)

// SyncUserDataHrsv syncs user data from HRSV
func (repo *repository) SyncUserDataHrsv(userData []dto.UserDataHRSV) error {
	emailMap := make(map[string]*entity.User)

	fmt.Println(userData)	

	timeNow := time.Now()

	// Convert dto.UserDataHRSV to entity.User
	for _, data := range userData {
		var role string

		if data.SystemRole == "User" {
			role = "user"
		}

		if data.CompanyRole.Name == "Project Manager" {
			role = "project_manager"
		}

		user := entity.User{
			Username:  data.Name,
			Phone:     "", // Set phone if available
			Workpalce: data.Workplace.Name,
			Team:      data.Team.Name,
			Email:     data.Email,
			Role:      role,
			CreatedBy: "system",
			UpdatedBy: "system",
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		}

		emailMap[user.Email] = &user
	}

	// Fetch existing users based on email
	var existingUsers []entity.User
	if err := repo.db.Where("email IN ?", getEmails(emailMap)).Find(&existingUsers).Error; err != nil {
		return err
	}

	// Update existing users and prepare for bulk insert
	for i := range existingUsers { // Use index to avoid copying
		existingUser := &existingUsers[i] // Get pointer to existing user
		if user, ok := emailMap[existingUser.Email]; ok {
			// Update existing user
			existingUser.Username = user.Username
			existingUser.Phone = user.Phone
			existingUser.Workpalce = user.Workpalce
			existingUser.Team = user.Team
			existingUser.Role = user.Role
			existingUser.UpdatedAt = &timeNow

			if err := repo.db.Save(existingUser).Error; err != nil {
				return err // Return error if save fails
			}
			delete(emailMap, existingUser.Email) // Remove from the map
		}
	}

	// Remaining users in emailMap are new users to insert
	newUsers := make([]entity.User, 0, len(emailMap))
	for _, user := range emailMap {
		newUsers = append(newUsers, *user)
	}

	// Perform bulk insert
	if len(newUsers) > 0 {
		if err := repo.db.Create(&newUsers).Error; err != nil {
			return err
		}
	}

	return nil
}

// Helper function to extract emails from user map
func getEmails(emailMap map[string]*entity.User) []string {
	emails := make([]string, 0, len(emailMap))
	for email := range emailMap {
		emails = append(emails, email)
	}
	return emails
}
