package helper

import (
	"os"

	"github.com/coderudu/cotach/database"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type signedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	UID        string
	User_type  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, user)

var secret_key string =os.Getenv("secretkey")

func GenerateAllTokens() {

}