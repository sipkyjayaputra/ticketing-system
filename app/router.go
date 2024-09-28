package app

import (
	"github.com/sipkyjayaputra/ticketing-system/delivery"
	"github.com/sipkyjayaputra/ticketing-system/middleware"
	"github.com/sipkyjayaputra/ticketing-system/repository"
	"github.com/sipkyjayaputra/ticketing-system/usecase"
	"github.com/sipkyjayaputra/ticketing-system/utils"
	"net/http"

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
		protectedRoutes.GET("users", middleware.AdminAccess(), del.GetUsers)
		protectedRoutes.POST("users", middleware.AdminAccess(), del.AddUser)
		protectedRoutes.PUT("users/:id", del.UpdateUser)
		protectedRoutes.GET("users/:id", del.GetUserById)
		protectedRoutes.DELETE("users/:id", middleware.AdminAccess(), del.DeleteUser)
	}
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
