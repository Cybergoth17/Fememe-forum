package routes

import (
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"test/controllers"
	_ "test/db"

	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/create", controllers.UserInsert())
	incomingRoutes.GET("/create", func(context *gin.Context) {
		context.HTML(http.StatusOK, "inserting.html", gin.H{})
	})
	incomingRoutes.GET("/read", controllers.ReadUser())
	incomingRoutes.GET("/read/:username", controllers.DeleteUser())
	incomingRoutes.GET("/readid/:username/:newusername", controllers.UpdateUser())
	incomingRoutes.GET("/readid/:username", controllers.ReadUserById())

	incomingRoutes.GET("/api/posts", controllers.SeeAllPosts())
	incomingRoutes.GET("/api/posts/:id", controllers.SeeSinglePost())
	incomingRoutes.GET("/api/posts/user/:username", controllers.SeePostsByUsername())
	incomingRoutes.POST("/api/posts", controllers.CreateUserPost())
	incomingRoutes.DELETE("/api/posts/:id", controllers.DeletePost())
	incomingRoutes.PUT("/api/posts/:id", controllers.UpdatePost())

	incomingRoutes.POST("/users/signup", controllers.Signup())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.GET("/users/signup", func(context *gin.Context) {
		context.HTML(http.StatusOK, "testRegister.html", gin.H{})
	})
	incomingRoutes.GET("/users/login", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", gin.H{})
	})
	incomingRoutes.GET("/users/signout", controllers.Logout())
}
