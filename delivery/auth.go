package delivery

import (
	"fmt"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-gonic/gin"
)

func (del *delivery) SignIn(c *gin.Context) {
	funcName := "SignIn"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	req := dto.SignIn{}
	err := c.ShouldBindJSON(&req)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", dto.SignIn{Email: req.Email, Password: "******"}), del.logger)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", err.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", err.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	res, errRes := del.uc.SignIn(req)

	if errRes != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process failed %s", errRes.Response.Errors), del.logger)
		c.JSON(errRes.Response.StatusCode, errRes)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) RefreshToken(c *gin.Context) {
	funcName := "RefreshToken"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	request := dto.RefreshToken{}
	errBind := c.ShouldBindJSON(&request)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", request), del.logger)

	if errBind != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", errBind.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errBind.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	res, err := del.uc.RefreshToken(request.RefreshToken)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}
