package configuration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func ConnectMySQL() (*gorm.DB, *sql.DB, error) {
	dsn := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`,
		CONFIG["MYSQL_USER"],
		CONFIG["MYSQL_PASS"],
		CONFIG["MYSQL_HOST"],
		CONFIG["MYSQL_PORT"],
		CONFIG["MYSQL_DB_NAME"],
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	sqlConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalln("Error connect to MySQL: ", err.Error())
		return nil, nil, err
	}
	sqlConn.AutoMigrate(&entity.User{})     // Migrate User first
	sqlConn.AutoMigrate(&entity.Document{}) // Migrate Ticket last
	sqlConn.AutoMigrate(&entity.Activity{}) // Migrate Document second
	sqlConn.AutoMigrate(&entity.Ticket{})   // Migrate Activity before Ticket

	sqlDb, errDb := sqlConn.DB()
	if errDb != nil {
		log.Println(errDb)
	} else {
		sqlDb.SetConnMaxIdleTime(300000)
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(1000)
	}

	log.Println("MySQL connection success")

	return sqlConn, sqlDb, nil
}
