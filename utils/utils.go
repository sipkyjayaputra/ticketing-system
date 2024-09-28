package utils

import (
	"fmt"
	"strings"

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
