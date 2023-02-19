package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/coderudu/cotach/Database"
	helper "github.com/coderudu/cotach/Helper"
	models "github.com/coderudu/cotach/Models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()

func VerifyPassword(userPassword, foundUserPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(foundUserPassword))

	check := true

	msg := ""

	if err != nil {

		msg = fmt.Sprint("email of password is incorrect")
		check = false
	}

	return check, msg

}

func HashPassword(password string) string{

	bytes, err:= bcrypt.GenerateFromPassword([]byte(password), 14)

	if err !=nil {
		log.Panic(err)
	}
	return string(bytes)

}

func Signup() gin.HandlerFunc{

	return func(c *gin.Context) {
		
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		if err := c.BindJSON(&user); err!=nil{
              
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr !=nil {
               c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		
		defer cancel()

		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to retrieve email of user"})
		}

		password :=HashPassword(*user.Password)

		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Password})

		defer cancel()
		if err !=nil{
        log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retireve phone number of the user" })
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This phone number has already been used before"})
		}

		user.CreatedAt, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339) )

		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserID, USER.UserType)

		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			msg := fmt.Sprintf("user item was not created")

			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)




	}


}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"email": "email or password is incorrect"})
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		token, refreshToken, err := helper.GenerateAllTokens(*foundUser.FirstName, *foundUser.LastName, *foundUser.Email, foundUser.UserID, *foundUser.UserType)
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserID}).Decode(foundUser.UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)

	}

}



func GetUserByID() {

}

func GetUsers() {

}
