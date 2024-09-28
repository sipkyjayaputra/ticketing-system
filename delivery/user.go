package delivery

import (
	"fmt"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"
	"time"

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
