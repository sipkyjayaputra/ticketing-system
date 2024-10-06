package delivery

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-gonic/gin"
)

func (del *delivery) GetActivitiesByTicketNo(c *gin.Context) {
	funcName := "GetActivitiesByTicketNo"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	ticketID := c.Param("ticket_id")
	res, err := del.uc.GetActivitiesByTicketNo(ticketID)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) AddActivity(c *gin.Context) {
	funcName := "AddActivity"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	form, errForm := c.MultipartForm()
	if errForm != nil {
		utils.LoggerProcess("error",  errForm.Error(), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errForm.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	documents := []dto.Document{}
	for key, val := range form.File {
		for _, f := range val {
			documents = append(documents, dto.Document{
				DocumentType: key,
				DocumentFile: f,
			})
		}
	}

	formValue := form.Value
	
	ticketID, _ := strconv.ParseUint(formValue["ticket_id"][0], 10, 64)
	newActivity := dto.Activity{
		TicketID:   uint(ticketID),
		Description: formValue["description"][0],
		Documents:   documents,
	}

	creator, _ := c.Get("user_id")
	res, err := del.uc.AddActivity(newActivity, creator.(uint))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) UpdateActivity(c *gin.Context) {
	funcName := "UpdateActivity"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	form, errForm := c.MultipartForm()
	if errForm != nil {
		utils.LoggerProcess("error", errForm.Error(), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errForm.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	documents := []dto.Document{}
	for key, val := range form.File {
		for _, f := range val {
			documents = append(documents, dto.Document{
				DocumentType: key,
				DocumentFile: f,
			})
		}
	}

	formValue := form.Value
	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ticketID, _ := strconv.ParseUint(formValue["ticket_id"][0], 10, 64)
	newActivity := dto.Activity{
		TicketID:    uint(ticketID),
		ActivityID:  uint(activityID),
		Description: formValue["description"][0],
		Documents:   documents,
	}

	updater, _ := c.Get("user_id")
	res, err := del.uc.UpdateActivity(newActivity, updater.(uint), uint(activityID))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) DeleteActivity(c *gin.Context) {
	funcName := "DeleteActivity"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	res, err := del.uc.DeleteActivity(uint(activityID))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) GetActivityById(c *gin.Context) {
	funcName := "GetActivityById"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	activityID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	res, err := del.uc.GetActivityById(uint(activityID))

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}
