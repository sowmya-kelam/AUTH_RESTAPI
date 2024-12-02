package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	"restapi/internal/controller"
	"restapi/internal/service"

	_ "restapi/docs"

	_ "github.com/mattn/go-sqlite3"

)

// Registering Routes 
func RegisterRoutes(server *gin.Engine, database *sql.DB) {

	service := service.NewUserService(database)
	userController := &controller.UserController{Service:service}

	userRoutes := server.Group("/")
    {
        userRoutes.POST("/signup", userController.Signup)
        userRoutes.POST("/login", userController.Login)
        userRoutes.POST("/authorize-token", userController.AuthorizeToken)
        userRoutes.POST("/revoke-token", userController.RevokeToken)
        userRoutes.POST("/refresh-token", userController.RefreshToken)
    }
	// Swagger API endpoint
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}