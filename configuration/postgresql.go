package configuration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/model/entity"

	"gorm.io/driver/postgres" // Import PostgreSQL driver
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres() (*gorm.DB, *sql.DB, error) {
	dsn := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta`,
		CONFIG["POSTGRES_HOST"],
		CONFIG["POSTGRES_PORT"],
		CONFIG["POSTGRES_USER"],
		CONFIG["POSTGRES_PASS"],
		CONFIG["POSTGRES_DB_NAME"],
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

	sqlConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger}) // Change to postgres driver
	if err != nil {
		log.Fatalln("Error connecting to PostgreSQL: ", err.Error())
		return nil, nil, err
	}
	sqlConn.AutoMigrate(&entity.User{})     // Migrate User first
	sqlConn.AutoMigrate(&entity.Document{}) // Migrate Document second
	sqlConn.AutoMigrate(&entity.Activity{}) // Migrate Activity third
	sqlConn.AutoMigrate(&entity.Ticket{})   // Migrate Ticket last

	sqlDb, errDb := sqlConn.DB()
	if errDb != nil {
		log.Println(errDb)
	} else {
		sqlDb.SetConnMaxIdleTime(300000)
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(1000)
	}

	log.Println("PostgreSQL connection success")

	return sqlConn, sqlDb, nil
}
