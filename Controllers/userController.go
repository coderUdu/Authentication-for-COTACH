package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coderudu/cotach/Helper"
	"github.com/coderudu/cotach/Models"
	"github.com/coderudu/cotach/Database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")


func VerifyPassword(userPassword, foundUserPassword string) (bool, string){

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(foundUserPassword))

	check := true 

	msg := ""

	if err != nil {
		
		msg = fmt.Sprint("email of password is incorrect")
		check = false
	}

	return check, msg

}

func PasswordHash(){

}

func Login() gin.HandlerFunc{

	return func(c *gin.Context){
		
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user) ;  err!=nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		
		defer cancel()

		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"email":"email or password is incorrect"})
		}

		passwordIsValid, msg:= VerifyPassword(*user.Password, *foundUser.Password)

		if  passwordIsValid !=true {
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

		if err!=nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
          c.JSON(http.StatusOK, foundUser)
		
	}


}

func Signup(){

}

func GetUserByID(){

}

func GetUsers(){
	
}

