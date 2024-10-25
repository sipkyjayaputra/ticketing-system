package usecase

import (
	"time"

	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) SignIn(request dto.SignIn) (*utils.ResponseContainer, *utils.ErrorContainer) {
	user, err := uc.repo.SignIn(request.Email)

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to sign in", err.Error())
	}

	return uc.handleDefaultSignIn(request.Password, user)
}

func (uc *usecase) RefreshToken(token string) (*utils.ResponseContainer, *utils.ErrorContainer) {
	var user entity.User
	var err error

	decodedUser, err := helpers.DecodeJWTToken(token)
	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to refresh token", err.Error())
	}

	user, err = uc.repo.SignIn(decodedUser.Username)

	if err != nil {
		return nil, utils.BuildInternalErrorResponse("failed to get user data", err.Error())
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

func (uc *usecase) handleDefaultSignIn(password string, user entity.User) (*utils.ResponseContainer, *utils.ErrorContainer) {
	if user.Email == "" {
		return nil, utils.BuildUnauthorizedResponse("invalid email", "user not found")
	}

	if user.Password == "" {
		return nil, utils.BuildUnauthorizedResponse("password not synced", "user not synced")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.BuildUnauthorizedResponse("invalid password", err.Error())
	}

	if user.Role == "" || user.Role == "user" {
		return nil, utils.BuildUnauthorizedResponse("invalid access", "user not allowed to access this application")
	}

	return uc.generateTokens(user)
}

func (uc *usecase) generateTokens(user entity.User) (*utils.ResponseContainer, *utils.ErrorContainer) {
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
