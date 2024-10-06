package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerProcess(logType string, message string, logger *logrus.Logger) {
	formattedMessage := fmt.Sprintf("[%s] %s", strings.ToUpper(logType), message)

	switch logType {
	case "warn":
		logger.Warn(formattedMessage)
	case "error":
		logger.Error(formattedMessage)
	case "info":
		logger.Info(formattedMessage)
	default:
		logger.Info("Unknown log type: " + logType)
	}
}


var mapCode = map[string]string{
	"Pengesahan Dokumen":   "DOC",
	"Notulensi Meeting":     "NOT",
	"Pelaporan Insiden":     "INC",
	"Pengelolaan Aset":      "INV",
	"Complaint Handling":     "COM",
	"Laporan Mutu Bulanan":  "ADD",
	"Request Fulfillment":    "REQ",
	"Pengajuan RFC":         "RFC",
}

// GenerateTicketNumber generates a ticket number based on the provided parameters.
func GenerateTicketNumber(idSerial uint, ticketType string) string {
	code, exists := mapCode[ticketType]
	if !exists {
		return "Unknown ticket type"
	}

	// Get current date
	currentTime := time.Now()
	monthRoman := getRomanMonth(currentTime.Month())
	year := currentTime.Year()

	// Format the ticket number
	ticketNumber := fmt.Sprintf("%d/SV-HR/%s/%s/%d", idSerial, code, monthRoman, year)
	return ticketNumber
}
func GenerateDocNumber(idTicket uint, ticketType string, order int) string {
	code, exists := mapCode[ticketType]
	if !exists {
		return "Unknown ticket type"
	}

	// Get current date
	currentTime := time.Now()
	monthRoman := getRomanMonth(currentTime.Month())
	year := currentTime.Year()

	// Format the ticket number
	ticketNumber := fmt.Sprintf("%d/SV-HR/%d-%s/%s/%d", idTicket, order, code, monthRoman, year)
	return ticketNumber
}

// getRomanMonth converts a month to its Roman numeral representation.
func getRomanMonth(month time.Month) string {
	months := []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	return months[month-1] // months are 1-based in time.Month
}