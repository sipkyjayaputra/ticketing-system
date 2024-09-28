package repository

import (
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
)

func (repo *repository) GetUsers() ([]entity.User, error) {
	users := []entity.User{}
	if err := repo.db.Model(&entity.User{}).Omit("Password").Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *repository) AddUser(user dto.User) error {
	newUser := &entity.User{
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
	}

	if err := repo.db.Create(&newUser).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) UpdateUser(user dto.User, id string) error {
	updateUser := &entity.User{
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		UpdatedBy: user.UpdatedBy,
	}

	if err := repo.db.Model(&entity.User{}).Where("id = ?", id).Updates(&updateUser).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) DeleteUser(id string) error {
	if err := repo.db.Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) GetUserById(id string) (*entity.User, error) {
	user := &entity.User{}
	if err := repo.db.Model(&entity.User{}).Where("id = ?", id).Omit("Password").First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
