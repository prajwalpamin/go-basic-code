package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/controller"
	_ "gitlab.niveussolutions.com/prajwal.amin/gop1/docs"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/middlewares"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/home", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "this is home page"})
		})
		api.GET("/users", controller.Users)
		api.GET("/users/:id", controller.User)
		api.PUT("/users/:id", controller.Update)
		api.POST("/register", controller.RegisterUser)
		api.DELETE("/users/:id", controller.DeleteUser)

		api.POST("/token", controller.GenrateToken)

		secured := api.Group("/secure").Use(middlewares.Auth())
		{
			secured.GET("/page", controller.SecuredPage)
			api.POST("/login", controller.LoginUser)

		}

	}
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
