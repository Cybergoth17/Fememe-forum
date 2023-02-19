package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test/db"
	database "test/db"
	"test/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

type CreatePostBody struct {
	Username string
	Title    string
	Text     string
	Tags     []string
	Avatar   string
}

type CreateCommentBody struct {
	Username string
	Text     string
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
		a, _ := c.Params.Get("id")
		oid, _ := primitive.ObjectIDFromHex(a)
		x, _ := FindOnePost(ctx, oid)
		c.JSON(200, x)
		defer cancel()
	}
}

func SeePostsByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		a, _ := c.Params.Get("username")
		log.Println(a)

		x, _ := FindPostsByUsername(ctx, a)
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

		log.Println(requestBody)

		post := models.Post{
			Username: requestBody.Username,
			Avatar:   requestBody.Avatar,
			Title:    requestBody.Title,
			Text:     requestBody.Text,
			Comment:  g,
			Tags:     requestBody.Tags,
			Date:     time.Now(),
		}
		_, insertErr := postCollection.InsertOne(ctx, post)
		if insertErr != nil {
			msg := fmt.Sprintf("Post item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		resultText := fmt.Sprintf("successfully created a post")
		c.JSON(http.StatusOK, gin.H{"message": resultText})
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
		c.JSON(http.StatusOK, gin.H{"message": resultText})

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
		c.JSON(http.StatusOK, gin.H{"message": resultText})
	}
}

func SeeAllComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var commentCollection *mongo.Collection = db.OpenCollection(db.Client, "comments")
		result, _ := commentCollection.Find(ctx, bson.M{})
		var u []models.Comment

		err := result.All(ctx, &u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find comments"})
		}
		c.JSON(200, u)
		defer cancel()
	}
}

func SeeAllCommentsByPostId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var commentCollection *mongo.Collection = db.OpenCollection(db.Client, "comments")
		id, _ := c.Params.Get("id")
		oid, _ := primitive.ObjectIDFromHex(id)
		result, _ := commentCollection.Find(ctx, bson.M{"post_id": oid})
		var u []models.Comment

		err := result.All(ctx, &u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find comments"})
		}
		c.JSON(200, u)
		defer cancel()
	}
}

func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentCollection *mongo.Collection = db.OpenCollection(db.Client, "comments")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		id, _ := c.Params.Get("id")
		oid, _ := primitive.ObjectIDFromHex(id)
		deleteRes, err := commentCollection.DeleteOne(ctx, bson.M{"_id": oid})
		if err != nil {
			msg := fmt.Sprintf("User post doesn't exist")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		resultText := fmt.Sprintf("successfully deleted %v post", deleteRes.DeletedCount)
		c.JSON(http.StatusOK, gin.H{"message": resultText})

	}
}

func CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var commentCollection *mongo.Collection = db.OpenCollection(db.Client, "comments")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var requestBody CreateCommentBody

		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println(err)
			return
		}

		id, _ := c.Params.Get("id")
		oid, _ := primitive.ObjectIDFromHex(id)
		post, err := FindOnePost(ctx, oid)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		comment := models.Comment{
			ID:       primitive.NewObjectID(),
			Text:     requestBody.Text,
			Date:     time.Now(),
			Username: requestBody.Username,
			PostID:   post.ID,
		}

		commentInsertRes, err := commentCollection.InsertOne(ctx, comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		log.Println("Inserted ID", commentInsertRes.InsertedID)

		var createdComment models.Comment
		err = commentCollection.FindOne(ctx, bson.M{"_id": commentInsertRes.InsertedID}).Decode(&createdComment)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		newCommentArr := append(post.Comment, createdComment)
		result := postCollection.FindOneAndUpdate(ctx, bson.M{"_id": post.ID}, bson.D{{"$set", bson.D{{"comments", newCommentArr}}}})

		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, gin.H{"message": "successfull created comment"})
	}
}
