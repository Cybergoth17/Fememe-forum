package routes

import (
	"net/http"
	"test/controllers"
	_ "test/db"

	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"

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

	incomingRoutes.GET("/api/posts/title/:title", controllers.SeePostsByTitle())

	incomingRoutes.POST("/api/posts/comment/:id", controllers.CreateComment())
	incomingRoutes.GET("/api/posts/comment/:id", controllers.SeeAllCommentsByPostId())
	incomingRoutes.GET("/api/posts/comment/", controllers.SeeAllComments())
	incomingRoutes.DELETE("/api/posts/comment/:id", controllers.DeleteComment())

	incomingRoutes.POST("/users/signup", controllers.Signup())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.GET("/users/signout", controllers.Logout())
	incomingRoutes.GET("/users", controllers.SeeAllUsers())
}
