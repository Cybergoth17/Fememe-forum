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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = db.OpenCollection(db.Client, "user")

func UserInsert() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var s = c.PostForm("username")
		var a = c.PostForm("password")
		var x = c.PostForm("email")
		user := models.User{
			Username: &s,
			Password: &a,
			Email:    &x,
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this username already exists"})
			return
		}

		user.ID = primitive.NewObjectID()
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}

}

func FindUser(ctx context.Context) (u []models.User, e error) {
	result, err := userCollection.Find(ctx, bson.M{})
	if result.Err() != nil {
		return u, fmt.Errorf("failed to find posts")
	}

	err = result.All(ctx, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func UserUpdate(ctx context.Context, username string, newusername string) (u models.User, e error) {
	result := userCollection.FindOneAndUpdate(ctx, bson.M{"username": username}, bson.D{{"$set", bson.D{{"username", newusername}}}})
	fmt.Println(result)
	return u, nil
}

func FindOne(ctx context.Context, username string) (u models.User, e error) {
	result := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	fmt.Println(result)
	return u, nil
}

func FindAllPost(ctx context.Context) (u []models.Post, e error) {
	var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	result, err := postCollection.Find(ctx, bson.M{}, opts)
	if result.Err() != nil {
		return u, fmt.Errorf("failed to find posts")
	}

	err = result.All(ctx, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func FindOnePost(ctx context.Context, id primitive.ObjectID) (u models.Post, e error) {
	var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")
	err := postCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&u)

	if err != nil {
		return u, err
	}

	return u, nil
}

func FindPostsByUsername(ctx context.Context, username string) (u []models.Post, e error) {
	var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	result, err := postCollection.Find(ctx, bson.M{"username": username}, opts)
	if err != nil {
		return u, fmt.Errorf("failed to find posts")
	}

	err = result.All(ctx, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func FindPostsByTitle(ctx context.Context, title string) (u []models.Post, e error) {
	var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	query := bson.M{"title": bson.M{"$regex": title, "$options": "im"}}
	result, err := postCollection.Find(ctx, query, opts)
	if err != nil {
		return u, fmt.Errorf("failed to find posts")
	}

	err = result.All(ctx, &u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteById(ctx context.Context, id string) (a int64, e error) {
	var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("primitive.ObjectIDFromHex ERROR:", err)
		return 0, err
	}
	res, err := postCollection.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return res.DeletedCount, nil
}

func UpdateById(ctx context.Context, id string, text string) (a int64, e error) {
	var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("primitive.ObjectIDFromHex ERROR:", err)
		return 0, err
	}
	res, err := postCollection.UpdateOne(ctx, bson.M{"_id": idPrimitive},
		bson.D{{Key: "$set", Value: bson.D{{Key: "text", Value: text}}}})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return res.ModifiedCount, nil
}
