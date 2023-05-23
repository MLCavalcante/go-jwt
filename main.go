package main

import (
	"github.com/MLCavalcante/go-jwt/controllers"
	"github.com/MLCavalcante/go-jwt/initializers"
	"github.com/MLCavalcante/go-jwt/middleware"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()

}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.Run() 
}