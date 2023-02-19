package main

import (
	"os"
	"test/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
	}))
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/images", "./templates/images/")
	routes.UserRoutes(router)
	router.Run(":" + port)
}
