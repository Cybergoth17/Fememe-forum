package controllers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"test/helpers"
	"test/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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
		resultInsertionNumber, insertErr := userCollection.DeleteOne(ctx, bson.M{"username": a})
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		fmt.Println(resultInsertionNumber)
		c.Redirect(303, "/read")
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
		var s = c.PostForm("username")
		var a = c.PostForm("password")
		var x = c.PostForm("email")
		user := models.User{
			Username: &s,
			Password: &a,
			Email:    &x,
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

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.SetCookie("username", *user.Username, 3600, "/", "localhost", false, false)
		fmt.Println(resultInsertionNumber)

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
		var s = c.PostForm("email")
		var a = c.PostForm("password")
		user := models.User{
			Email:    &s,
			Password: &a,
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
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.Username)

		helpers.UpdateAllTokens(token, refreshToken, *foundUser.Email)
		c.SetCookie("token", refreshToken, 3600, "/", "localhost", false, false)
		c.SetCookie("username", *foundUser.Username, 3600, "/", "localhost", false, false)

		c.Redirect(303, "/read")

	}
}
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("username", "expired", -1, "/", "localhost", false, false)
		c.SetCookie("token", "expiredToken", -1, "/", "localhost", false, false)
		c.Redirect(http.StatusFound, "/read")
	}
}
