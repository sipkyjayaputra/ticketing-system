package usecase

import (
	"strconv"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) SyncUserDataHrsv(userData []dto.UserDataHRSV) (*utils.ResponseContainer, *utils.ErrorContainer) {
	err := uc.repo.SyncUserDataHrsv(userData)

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to sync users", err.Error())
	}


	return utils.BuildSuccessResponse(err), nil
}

func (uc *usecase) SyncPasswordHrsv(email, password string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	// Mendapatkan informasi pengguna berdasarkan email
	user, err := uc.repo.SignIn(email)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get user info", err.Error())
	}

	// Menghasilkan hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.BuildBadRequestResponse("failed to generate password", err.Error())
	}


	userIDString := strconv.FormatUint(uint64(user.ID), 10)
	// Menyusun request untuk memperbarui password
	request := dto.UpdateUserPassword{
		ID:          userIDString,
		NewPassword: string(hashedPassword), // Menggunakan hashed password langsung
	}

	// Memperbarui password pengguna
	errUpdate := uc.repo.UpdateUserPassword(request)
	if errUpdate != nil {
		return nil, utils.BuildInternalErrorResponse("failed to update user password", errUpdate.Error())
	}

	if user.Role == "" || user.Role == "user" {
		return nil, utils.BuildUnauthorizedResponse("invalid access", "user not allowed to access this application")
	}

	accessToken, err := helpers.GenerateJWTToken(user, time.Hour*1)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to generate access token", err.Error())
	}

	refreshToken, err := helpers.GenerateJWTToken(user, time.Hour*3)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to generate refresh token", err.Error())
	}

	return utils.BuildSuccessResponse(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}), nil
}
