package delivery

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-gonic/gin"
)

func (del *delivery) GetUsers(c *gin.Context) {
	funcName := "GetUsers"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	res, err := del.uc.GetUsers()

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) AddUser(c *gin.Context) {
	funcName := "AddUser"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	request := dto.User{}
	errBind := c.ShouldBindJSON(&request)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", request), del.logger)

	if errBind != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", errBind.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errBind.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	if request.Password == "" {
		utils.LoggerProcess("error", "Invalid cannot be empty", del.logger)
		resp := utils.BuildBadRequestResponse("bad request", "password required")
		c.JSON(resp.Response.StatusCode, resp)
	}

	res, err := del.uc.AddUser(request)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) UpdateUser(c *gin.Context) {
	funcName := "UpdateUser"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	request := dto.User{}
	errBind := c.ShouldBindJSON(&request)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", request), del.logger)

	if errBind != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", errBind.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errBind.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	updater, _ := c.Get("username")
	id := c.Param("id")

	res, err := del.uc.UpdateUser(request, updater.(string), id)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) DeleteUser(c *gin.Context) {
	funcName := "DeleteUser"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	id := c.Param("id")

	res, err := del.uc.DeleteUser(id)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) GetUserById(c *gin.Context) {
	funcName := "GetUserById"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	id := c.Param("id")

	res, err := del.uc.GetUserById(id)
	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) UpdateUserPhoto(c *gin.Context) {
	funcName := "UpdateUserPhoto"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	form, errForm := c.MultipartForm()
	if errForm != nil {
		utils.LoggerProcess("error", fmt.Sprintf("%s", errForm.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errForm.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	photo := form.File["photo"]
	if len(photo) == 0 {
		errPhoto := errors.New("invalid photo")
		utils.LoggerProcess("error", fmt.Sprintf("%s", errPhoto.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errPhoto.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	id := c.Param("id")
	userID, errUserId := strconv.ParseInt(id, 10, 64)
	if errUserId != nil {
		utils.LoggerProcess("error", fmt.Sprintf("%s", errUserId.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errUserId.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	claimId, _ := c.Get("user_id")
	claimRole, _ := c.Get("role")
	if !strings.Contains(claimRole.(string), "admin") && claimId.(string) != id {
		err := errors.New("unauthorized access")
		utils.LoggerProcess("error", fmt.Sprintf("%s", err.Error()), del.logger)
		resp := utils.BuildForbiddenAccessResponse(err.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	userPhoto := dto.UpdateUserPhoto{
		ID:    uint(userID),
		Photo: photo[0],
	}
	res, err := del.uc.UpdateUserPhoto(userPhoto)
	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) UpdateUserPassword(c *gin.Context) {
	funcName := "UpdateUserPassword"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	request := dto.UpdateUserPassword{}
	errBind := c.ShouldBindJSON(&request)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", request), del.logger)
	if errBind != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", errBind.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errBind.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	id := c.Param("id")
	userID, errUserId := strconv.ParseInt(id, 10, 64)
	if errUserId != nil {
		utils.LoggerProcess("error", fmt.Sprintf("%s", errUserId.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errUserId.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	claimId, _ := c.Get("user_id")
	claimRole, _ := c.Get("role")
	if !strings.Contains(claimRole.(string), "admin") && claimId.(string) != id {
		err := errors.New("unauthorized access")
		utils.LoggerProcess("error", fmt.Sprintf("%s", err.Error()), del.logger)
		resp := utils.BuildForbiddenAccessResponse(err.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	if request.NewPassword != request.VerifyPassword {
		errPhoto := errors.New("password did not match")
		utils.LoggerProcess("error", fmt.Sprintf("%s", errPhoto.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errPhoto.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	request.ID = uint(userID)
	res, err := del.uc.UpdateUserPassword(request)
	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}
