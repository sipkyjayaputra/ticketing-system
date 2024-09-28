package repository

import "github.com/sipkyjayaputra/ticketing-system/model/entity"

func (repo *repository) SignIn(username string) (entity.User, error) {
	user := entity.User{}
	if err := repo.db.Where("username = ?", username).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
