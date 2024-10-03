package delivery

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cs "github.com/sipkyjayaputra/ticketing-system/constants"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (del *delivery) GetTickets(c *gin.Context) {
	funcName := "GetTickets"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	ticketFilter := dto.TicketFilter{
		TicketType:      c.Query("ticket_type"),
		Priority:        c.Query("priority"),
		Status:          c.Query("status"),
		ReportStartDate: c.Query("report_start_date"),
		ReportEndDate:   c.Query("report_end_date"),
		Terms:           c.Query("terms"),
		Limit:           c.Query("limit"),
		Offset:          c.Query("offset"),
	}

	res, err := del.uc.GetTickets(ticketFilter)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) GetTicketSummary(c *gin.Context) {
	funcName := "GetTicketSummary"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	res, err := del.uc.GetTicketSummary()

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

	form, errForm := c.MultipartForm()
	if errForm != nil {
		utils.LoggerProcess("error", fmt.Sprintf("%s", errForm.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errForm.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	id := uuid.New()
	formValue := form.Value
	reporterID, _ := strconv.ParseInt(formValue["reporter_id"][0], 10, 64)
	reportDate, _ := time.Parse(cs.DATE_TIME_LAYOUT, formValue["report_date"][0])
	assignedID, _ := strconv.ParseInt(formValue["assigned_id"][0], 10, 64)
	request := dto.Ticket{
		TicketNo:   id.String(),
		ReporterID: uint(reporterID),
		TicketType: formValue["ticket_type"][0],
		Subject:    formValue["subject"][0],
		ReportDate: reportDate,
		AssignedID: uint(assignedID),
		Priority:   formValue["priority"][0],
		Status:     formValue["status"][0],
		Content:    json.RawMessage(formValue["content"][0]),
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

	newActivity := dto.Activity{
		Documents: documents,
	}
	request.Activities = []dto.Activity{newActivity}

	creator, _ := c.Get("user_id")
	res, err := del.uc.AddTicket(request, creator.(uint))

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

	form, errForm := c.MultipartForm()
	if errForm != nil {
		utils.LoggerProcess("error", fmt.Sprintf("%s", errForm.Error()), del.logger)
		resp := utils.BuildBadRequestResponse("bad request", errForm.Error())
		c.JSON(resp.Response.StatusCode, resp)
		return
	}

	formValue := form.Value
	reporterID, _ := strconv.ParseInt(formValue["reporter_id"][0], 10, 64)
	reportDate, _ := time.Parse(cs.DATE_TIME_LAYOUT, formValue["report_date"][0])
	assignedID, _ := strconv.ParseInt(formValue["assigned_id"][0], 10, 64)
	request := dto.Ticket{
		ReporterID: uint(reporterID),
		TicketType: formValue["ticket_type"][0],
		Subject:    formValue["subject"][0],
		ReportDate: reportDate,
		AssignedID: uint(assignedID),
		Priority:   formValue["priority"][0],
		Status:     formValue["status"][0],
		Content:    json.RawMessage(formValue["content"][0]),
	}

	updater, _ := c.Get("user_id")
	ticketNo := c.Param("id")
	res, err := del.uc.UpdateTicket(request, updater.(uint), ticketNo)

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

	ticketNo := c.Param("id")
	res, err := del.uc.DeleteTicket(ticketNo)

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

	ticketNo := c.Param("id")
	res, err := del.uc.GetTicketById(ticketNo)

	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}
