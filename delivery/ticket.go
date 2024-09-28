package delivery

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-gonic/gin"
)

func (del *delivery) GetTickets(c *gin.Context) {
	funcName := "GetTickets"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	res, err := del.uc.GetTickets()

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) AddTicket(c *gin.Context) {
	funcName := "AddTicket"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	request := dto.Ticket{}
	errBind := c.ShouldBindJSON(&request)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", request), del.logger)

	if errBind != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", errBind.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errBind.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	creator, _ := c.Get("user_id")
	creatorID, _ := strconv.ParseInt(creator.(string), 10, 64)
	res, err := del.uc.AddTicket(request, uint(creatorID))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) UpdateTicket(c *gin.Context) {
	funcName := "UpdateTicket"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	request := dto.Ticket{}
	errBind := c.ShouldBindJSON(&request)
	utils.LoggerProcess("info", fmt.Sprintf("Request Body: %+v", request), del.logger)

	if errBind != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Validation failed %s", errBind.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errBind.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	updater, _ := c.Get("user_id")
	updaterID, _ := strconv.ParseInt(updater.(string), 10, 64)
	ticket := c.Param("ticket_no")
	ticketNo, _ := strconv.ParseInt(ticket, 10, 64)

	res, err := del.uc.UpdateTicket(request, uint(updaterID), uint(ticketNo))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) DeleteTicket(c *gin.Context) {
	funcName := "DeleteTicket"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	ticket := c.Param("ticket_no")
	ticketNo, _ := strconv.ParseInt(ticket, 10, 64)

	res, err := del.uc.DeleteTicket(uint(ticketNo))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) GetTicketById(c *gin.Context) {
	funcName := "GetTicketById"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	ticket := c.Param("ticket_no")
	ticketNo, _ := strconv.ParseInt(ticket, 10, 64)

	res, err := del.uc.GetTicketById(uint(ticketNo))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}
