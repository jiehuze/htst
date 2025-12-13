package routers

import (
	"htst/internal/app/controllers"
	v1 "htst/internal/app/controllers/v1"
	"htst/internal/app/services"

	"sync"

	"github.com/gin-gonic/gin"
)

var apiOnce sync.Once
var g *gin.Engine

func SetUp() *gin.Engine {
	apiOnce.Do(func() {
		g = gin.Default()

		// 跨域中间件
		// g.Use(corsMiddleware())

		mainGroup := g.Group("/htst")
		mainGroup.GET("/health", controllers.Health)

		userController := v1.NewUserController(services.IUser)
		mainGroup.POST("/users/login", userController.AuthUser)
		mainGroup.GET("/users/list", userController.GetUserList)
		mainGroup.PUT("/users/update", userController.UpdateUser)
		mainGroup.POST("/users/create", userController.CreateUser)
		mainGroup.DELETE("/users/:id", userController.DeleteUser)

		infoController := v1.NewInfoController(services.IInfo)
		mainGroup.POST("/info/add", infoController.CreateInfo)
		mainGroup.DELETE("/info/:id", infoController.DeleteInfo)
		mainGroup.GET("/info/list", infoController.GetInfoList)
		mainGroup.GET("/info/click/:id", infoController.Increment)
	})

	return g
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
