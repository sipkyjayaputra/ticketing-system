package app

import (
	"net/http"

	"github.com/sipkyjayaputra/ticketing-system/delivery"
	"github.com/sipkyjayaputra/ticketing-system/middleware"
	"github.com/sipkyjayaputra/ticketing-system/repository"
	"github.com/sipkyjayaputra/ticketing-system/usecase"
	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, logger *logrus.Logger, cache *redis.Client) *gin.Engine {
	repo := repository.NewRepository(db, logger, cache)
	uc := usecase.NewUsecase(repo, logger)
	del := delivery.NewDelivery(uc, logger)

	router := gin.New()
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/ping"),
		gin.Recovery(),
	)

	router.Use(getDefaultCORS())
	protectedRoutes := router.Group("/")
	protectedRoutes.Use(middleware.Authorization())
	{
		// USERS
		protectedRoutes.GET("users", middleware.AdminAccess(), del.GetUsers)
		protectedRoutes.POST("users", middleware.AdminAccess(), del.AddUser)
		protectedRoutes.PUT("users/:id", del.UpdateUser)
		protectedRoutes.GET("users/:id", del.GetUserById)
		protectedRoutes.DELETE("users/:id", middleware.AdminAccess(), del.DeleteUser)

		// TICKET
		protectedRoutes.GET("tickets", middleware.AdminAccess(), del.GetTickets)
		protectedRoutes.POST("tickets", middleware.AdminAccess(), del.AddTicket)
		protectedRoutes.PUT("tickets/:id", del.UpdateTicket)
		protectedRoutes.GET("tickets/:id", del.GetTicketById)
		protectedRoutes.DELETE("tickets/:id", middleware.AdminAccess(), del.DeleteTicket)

		// ACTIVITY
		protectedRoutes.GET("activities", del.GetActivitiesByTicketNo) // Retrieve activities by ticket number
		protectedRoutes.POST("activities", del.AddActivity)            // Add activity to a ticket
		protectedRoutes.PUT("activities/:id", del.UpdateActivity)      // Update an existing activity
		protectedRoutes.GET("activities/:id", del.GetActivityById)     // Get a specific activity by ID
		protectedRoutes.DELETE("activities/:id", del.DeleteActivity)   // Delete an activity
	}
	// AUTH
	router.POST("/auth/sign-in", del.SignIn)
	router.POST("/auth/sign-up", del.AddUser)
	router.POST("/auth/refresh-token", del.RefreshToken)
	router.GET("/ping", Ping)

	return router
}

func getDefaultCORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Content-Disposition", "Authorization", "X-Chub-Personal-Number", "X-Brisim-Token"}
	return cors.New(config)
}

func Ping(c *gin.Context) {
	response := utils.BuildSuccessResponse(map[string]string{"message": "pong"})
	c.JSON(http.StatusOK, response)
}
