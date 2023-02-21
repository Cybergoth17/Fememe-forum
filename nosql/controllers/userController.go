package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"test/helpers"
	"test/models"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateUserBody struct {
	Avatar   string
	Username string
	Email    string
	Password string
}

type LoginBody struct {
	Email    string
	Password string
}

func ReadUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//x, _ := FindUserPosts(context.Background(), a)
		x, _ := FindUser(context.Background())
		c.HTML(http.StatusOK, "listing.html", x)
	}
}
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		a, _ := c.Params.Get("username")
		result, _ := FindOne(ctx, a)

		_, err := postCollection.DeleteMany(ctx, bson.M{"username": result.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete posts"})
			return
		}

		_, err = commentCollection.DeleteMany(ctx, bson.M{"username": result.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete comments"})
			return
		}

		_, err = userCollection.DeleteOne(ctx, bson.M{"username": a})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

	}
}
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		a, _ := c.Params.Get("newusername")
		b, _ := c.Params.Get("username")
		lux, err := UserUpdate(context.Background(), b, a)
		if err != nil {
			return
		}
		c.Redirect(303, "/read")
		fmt.Println(lux)
	}
}
func ReadUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context
		username, _ := c.Params.Get("username")
		result, _ := FindOne(ctx, username)
		c.JSON(http.StatusOK, result)
	}
}

var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var requestBody CreateUserBody

		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println(err)
			return
		}

		user := models.User{
			Username: &requestBody.Username,
			Password: &requestBody.Password,
			Email:    &requestBody.Email,
			Avatar:   &requestBody.Avatar,
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		user.ID = primitive.NewObjectID()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.Username)
		user.Token = &token
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
	}

}
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var requestBody LoginBody

		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Println(err)
			return
		}

		user := models.User{
			Email:    &requestBody.Email,
			Password: &requestBody.Password,
		}
		var foundUser models.User

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.Username)

		helpers.UpdateAllTokens(token, refreshToken, *foundUser.Email)
		c.JSON(http.StatusOK, gin.H{"token": token, "user": foundUser})
	}
}
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("username", "expired", -1, "/", "localhost", false, false)
		c.SetCookie("token", "expiredToken", -1, "/", "localhost", false, false)
		c.Redirect(http.StatusFound, "/read")
	}
}

func SeeAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, _ := userCollection.Find(ctx, bson.M{})
		var u []models.User

		err := result.All(ctx, &u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find users"})
		}
		c.JSON(200, gin.H{"users": u})
		defer cancel()
	}
}
