package usecase

import (
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) GetUsers() (*utils.ResponseContainer, *utils.ErrorContainer) {
	user, err := uc.repo.GetUsers()

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get users", err.Error())
	}

	return utils.BuildSuccessResponse(user), nil
}

func (uc *usecase) AddUser(user dto.User) (*utils.ResponseContainer, *utils.ErrorContainer) {
	hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errHashedPassword != nil {
		return nil, utils.BuildBadRequestResponse("failed to generate password", errHashedPassword.Error())
	}
	user.Password = string(hashedPassword)
	user.CreatedBy = user.Username
	user.UpdatedBy = user.Username

	if err := uc.repo.AddUser(user); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to add user", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) UpdateUser(user dto.User, updater string, id string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	user.UpdatedBy = updater
	if err := uc.repo.UpdateUser(user, id); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to update user", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) DeleteUser(id string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	if err := uc.repo.DeleteUser(id); err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to update user", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) GetUserById(id string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	user, err := uc.repo.GetUserById(id)

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get user", err.Error())
	}

	return utils.BuildSuccessResponse(user), nil
}

func (uc *usecase) UpdateUserPhoto(request dto.UpdateUserPhoto) (*utils.ResponseContainer, *utils.ErrorContainer) {
	err := uc.repo.UpdateUserPhoto(request)

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get user", err.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}

func (uc *usecase) UpdateUserPassword(request dto.UpdateUserPassword) (*utils.ResponseContainer, *utils.ErrorContainer) {
	currentPassword, err := uc.repo.GetUserPassword(request.ID)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get user password", err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(request.CurrentPassword)); err != nil {
		return nil, utils.BuildBadRequestResponse("invalid password", err.Error())
	}
	hashedPassword, errHashedPassword := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if errHashedPassword != nil {
		return nil, utils.BuildBadRequestResponse("failed to generate password", errHashedPassword.Error())
	}

	request.NewPassword = string(hashedPassword)

	errUpdate := uc.repo.UpdateUserPassword(request)

	if errUpdate != nil {
		return nil, utils.BuildInternalErrorResponse("failed to update user", errUpdate.Error())
	}

	return utils.BuildSuccessResponse(nil), nil
}
