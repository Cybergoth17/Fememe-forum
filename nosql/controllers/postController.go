package controllers

import (
	"context"
	"fmt"
	"net/http"
	database "test/db"
	"test/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

type CreatePostBody struct {
	Username string
	Title    string
	Text     string
	Tags     []string
}

type UpdatePostBody struct {
	Text string `json:"text"`
}

func SeeAllPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		a, _ := FindAllPost(ctx)
		c.JSON(200, a)
		defer cancel()
	}
}

func SeeSinglePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		a, _ := c.Params.Get("title")
		fmt.Println(a)
		x, _ := FindOnePost(ctx, a)
		c.JSON(200, x)
		defer cancel()
	}
}

func CreateUserPost() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		g := make([]models.Comment, 0, 0)

		var requestBody CreatePostBody

		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println(err)
			return
		}

		post1 := models.Post{
			Username: requestBody.Username,
			Title:    requestBody.Title,
			Text:     requestBody.Text,
			Comment:  g,
			Tags:     requestBody.Tags,
		}
		_, insertErr := postCollection.InsertOne(ctx, post1)
		if insertErr != nil {
			msg := fmt.Sprintf("Post item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		resultText := fmt.Sprintf("successfully created a post")
		c.JSON(http.StatusInternalServerError, gin.H{"message": resultText})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id, _ := c.Params.Get("id")
		deletedCount, err := DeleteById(ctx, id)
		if err != nil {
			msg := fmt.Sprintf("User post doesn't exist")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		resultText := fmt.Sprintf("successfully deleted %v post", deletedCount)
		c.JSON(http.StatusInternalServerError, gin.H{"message": resultText})

	}
}

func UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var requestBody UpdatePostBody

		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println(err)
			return
		}

		id, _ := c.Params.Get("id")
		text := requestBody.Text
		modifiedCount, err := UpdateById(ctx, id, text)
		if err != nil {
			msg := fmt.Sprintf("User post doesn't exist")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		resultText := fmt.Sprintf("successfully updated %v post", modifiedCount)
		c.JSON(http.StatusInternalServerError, gin.H{"message": resultText})
	}
}