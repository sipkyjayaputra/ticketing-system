package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/app"
	"github.com/sipkyjayaputra/ticketing-system/configuration"
)

func main() {
	configuration.InitConfig()

	sqlConn, sqlDb, errSql := configuration.ConnectMySQL()
	if errSql != nil {
		log.Println("Error: PosgreSQL connection failed. ", errSql)
	}
	defer sqlDb.Close()

	logger, logFile, errLogger := configuration.InitLogger()
	if errLogger != nil {
		log.Println("Error: Failed init logger. ", errLogger)
	}
	defer logFile.Close()

	cache, errCache := configuration.ConnectRedis()
	if errCache != nil {
		log.Println("Error: Failed init logger. ", errLogger)
	}
	defer cache.Close()

	log.Println("Routes Initialized")
	router := app.InitRouter(sqlConn, logger, cache)

	port := configuration.Getenv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Println("Server Initialized")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
