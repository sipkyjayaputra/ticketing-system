package configuration

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func InitLogger() (*logrus.Logger, *os.File, error) {
	// Pastikan folder logs/go/log sudah ada, jika belum, buat folder tersebut
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.MkdirAll("logs", os.ModePerm)
		if err != nil {
			return nil, nil, err
		}
	}

	// Dapatkan tanggal saat ini dalam format "YYYY-MM-DD"
	currentDate := time.Now().Format("2006-01-02")

	// Buat nama file log dengan format yang diinginkan
	logFileName := "logs/go-service-log-" + currentDate + ".log"

	// Open or create log file with the generated filename
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// Inisialisasi logger
	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	return logger, logFile, nil
}
