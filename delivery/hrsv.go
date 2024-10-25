package delivery

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sipkyjayaputra/ticketing-system/model/dto"
	"github.com/sipkyjayaputra/ticketing-system/utils"
)

func (del *delivery) GetHrsvUsers(c *gin.Context) {
	del.handleHrsvRequest(c, "/users", "GetHrsvUsers")
}

func (del *delivery) GetHrsvRoles(c *gin.Context) {
	del.handleHrsvRequest(c, "/roles", "GetHrsvRoles")
}

func (del *delivery) SyncUserDataHrsv(c *gin.Context) {
	funcName := "SyncUserDataHrsv"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	endpoint := "/users"
	response, err := del.hrsvClient.Query(endpoint)
	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Lower %s, [ERROR]: %v", funcName, err), del.logger)
		c.JSON(http.StatusInternalServerError, utils.BuildInternalErrorResponse(fmt.Sprintf("failed to get hrsv %s data", endpoint), err.Error()))
		return
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		utils.LoggerProcess("error", fmt.Sprintf("Lower %s, [ERROR]: data field not found or not a valid array", funcName), del.logger)
		c.JSON(http.StatusInternalServerError, utils.BuildInternalErrorResponse("invalid response format", "data field not found or not a valid array"))
		return
	}

	var hrsvUsers []dto.UserDataHRSV

	for _, item := range data {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Handle error or log invalid item
		}

		var user dto.UserDataHRSV
		user.ID = itemMap["_id"].(string)
		user.Email = itemMap["email"].(string)
		user.Name = itemMap["name"].(string)
		user.FID = itemMap["fid"].(string)
		user.IsActive = itemMap["is_active"].(bool)
		user.IsSynced = itemMap["is_synced"].(bool)
		user.IsVerified = itemMap["is_verified"].(bool)
		user.SystemRole = itemMap["system_role"].(string)

		if role, ok := itemMap["company_role"].(map[string]interface{}); ok {
			user.CompanyRole.FID = role["fid"].(string)
			user.CompanyRole.Name = role["name"].(string)
		}

		if team, ok := itemMap["team"].(map[string]interface{}); ok {
			user.Team.FID = team["fid"].(string)
			user.Team.Name = team["name"].(string)
		}

		if workplace, ok := itemMap["workplace"].(map[string]interface{}); ok {
			user.Workplace.FID = workplace["fid"].(string)
			user.Workplace.Name = workplace["name"].(string)
			if location, ok := workplace["location"].(map[string]interface{}); ok {
				user.Workplace.Location.Lat = location["lat"].(float64)
				user.Workplace.Location.Lng = location["lng"].(float64)
			}
			user.Workplace.Radius = int(workplace["radius"].(float64)) // Assuming radius is a float64
		}

		// Parse time
		if createdAt, ok := itemMap["created_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
				user.CreatedAt = t
			}
		}
		if updatedAt, ok := itemMap["updated_at"].(string); ok {
			if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				user.UpdatedAt = t
			}
		}

		hrsvUsers = append(hrsvUsers, user)
	}

	res, errRes := del.uc.SyncUserDataHrsv(hrsvUsers)

	if errRes != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", errRes.Response.Errors), del.logger)
		c.JSON(errRes.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) SyncPasswordHrsv(c *gin.Context) {
	funcName := "SyncUserDataHrsv"
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

	payload := map[string]string{
		"email" : req.Email,
		"password" : req.Password,
	}

	errAuth := del.hrsvClient.AuthenticateWithPayload(payload)

	if errAuth != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Lower %s, [ERROR]: %v", funcName, errAuth), del.logger)
		c.JSON(http.StatusInternalServerError, utils.BuildInternalErrorResponse("failed to get hrsv authenticate user", errAuth.Error()))
		return
	}

	res, errRes := del.uc.SyncPasswordHrsv(req.Email, req.Password)

	if errRes != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", errRes.Response.Errors), del.logger)
		c.JSON(errRes.Response.StatusCode, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(res.Response.StatusCode, res)
}

func (del *delivery) handleHrsvRequest(c *gin.Context, endpoint string, funcName string) {
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	// Call the Query method on hrsvClient
	response, err := del.hrsvClient.Query(endpoint)
	if err != nil {
		// Handle the error accordingly
		utils.LoggerProcess("error", fmt.Sprintf("Lower %s, [ERROR]: %v", funcName, err), del.logger)
		c.JSON(http.StatusInternalServerError, utils.BuildInternalErrorResponse(fmt.Sprintf("failed to get hrsv %s data", endpoint), err.Error()))
		return
	}

	// Extract data safely
	data, ok := response["data"].([]interface{})
	if !ok {
		utils.LoggerProcess("error", fmt.Sprintf("Lower %s, [ERROR]: data field not found or not a valid array", funcName), del.logger)
		c.JSON(http.StatusInternalServerError, utils.BuildInternalErrorResponse("invalid response format", "data field not found or not a valid array"))
		return
	}

	// Log and return the response to the client
	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.JSON(http.StatusOK, utils.BuildSuccessResponse(data))
}
