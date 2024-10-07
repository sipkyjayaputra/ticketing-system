package repository

import "github.com/sipkyjayaputra/ticketing-system/model/entity"

func (repo *repository) SignIn(email string) (entity.User, error) {
	user := entity.User{}
	if err := repo.db.Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
