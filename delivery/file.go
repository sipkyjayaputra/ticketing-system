package delivery

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sipkyjayaputra/ticketing-system/utils"
)

func (del *delivery) FileServe(c *gin.Context) {
	funcName := "FileServe"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	documentId := c.Param("id")
	if documentId == "" {
		err := errors.New("invalid id")
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Error()), del.logger)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, _ := strconv.ParseInt(documentId, 10, 64)
	document, err := del.uc.GetDocumentById(uint(id))
	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	if _, err := filepath.Abs(document.DocumentPath); err != nil {
		err := errors.New("file not found")
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Error()), del.logger)
		c.JSON(http.StatusNotFound, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.File(document.DocumentPath)
}

func (del *delivery) FileDownload(c *gin.Context) {
	funcName := "FileDownload"
	startTime := time.Now()
	utils.LoggerProcess("info", fmt.Sprintf("Upper %s, [START]: Processing Request", funcName), del.logger)

	documentId := c.Param("id")
	if documentId == "" {
		err := errors.New("invalid id")
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Error()), del.logger)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, _ := strconv.ParseInt(documentId, 10, 64)
	document, err := del.uc.GetDocumentById(uint(id))
	if err != nil {
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Response.Errors), del.logger)
		c.JSON(err.Response.StatusCode, err)
		return
	}

	if _, err := filepath.Abs(document.DocumentPath); err != nil {
		err := errors.New("file not found")
		utils.LoggerProcess("error", fmt.Sprintf("Process Failed %s", err.Error()), del.logger)
		c.JSON(http.StatusNotFound, err)
		return
	}

	utils.LoggerProcess("info", fmt.Sprintf("Lower %s, [END]: Elapsed Time %v", funcName, time.Since(startTime)), del.logger)
	c.FileAttachment(document.DocumentPath, document.DocumentName)
}
