package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `BSON:"id"`
	FirstName    *string            `json:"firstname" validate: "Required, min=2 max=200"`
	LastName     *string            `json:"firstname" validate: "Required, min=2 max=200"`
	Password     *string            `json:"password" validate: "Required, min=8"`
	PasswordHash *string            `json:"passwordhash" validate: "Required, min=8"`
	Email        *string            `json:"Email" validate: "Required"`
	PhoneNumber  *string            `json:"phonenumber" validate: "Required, min=11"`
	Token        *string            `json: "token" validate: "required"`
	RefreshToken *string            `json: "refreshtoken" validate: "required"`
	UserType     *string            `json: "usertype" validate: "Required,  eq=Admin|eq=user"`
	UserID       *string            `json: "token" validate: "required"`
	CreatedAt    time.Time          `json: "createdat"`
	UpdatedAt    time.Time          `json: "updatedat"`
}
