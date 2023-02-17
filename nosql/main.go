package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"test/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/images", "./templates/images/")
	routes.UserRoutes(router)
	router.Run(":" + port)
}
