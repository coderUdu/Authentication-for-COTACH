package controllers

import (
	"context"
)



func VerifyPassword(){

}

func PasswordHash(){

}

func Login(){

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
			c.JSON(http.StatusInternalServerError, gin.H{"email": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
		}

		defer cancel()

		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		token, refreshtoken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_id, *foundUser.User_Type)
        helper.UpdateAllTokens( token, refreshtoken, foundUser.User_id)
        err = userCollection.FindOne(ctx, bson.M{ "user_id": foundUser.User_id}).Decode(&foundUser)

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

