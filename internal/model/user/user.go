package modelUser

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type LoginUser struct {
// 	ID       primitive.ObjectID `bson:"_id"`
// 	Password *string            `json:"password" validate:"required,min=6,max=100"`
// 	Email    *string            `json:"email" validate:"email,required"`
// }

// type SignupUser struct {
// 	ID           primitive.ObjectID `bson:"_id"`
// 	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
// 	Password     *string            `json:"password" validate:"required,min=6,max=100"`
// 	Email        *string            `json:"email" validate:"email,required"`
// 	Avator       *string            `json:"avator"`
// 	AccessToken  *string            `json:"accessToken"`
// 	RefreshToken *string            `json:"refreshToken"`
// 	UserType     *string            `json:"userType" validate:"required,eq=ADMIN|eq=USER"`
// 	CreatedAt    time.Time          `json:"createdAt"`
// 	UpdatedAt    time.Time          `json:"updatedAt"`
// 	UserId       string             `json:"userId"`
// }
// type User struct {
// 	ID           primitive.ObjectID `bson:"_id"`
// 	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
// 	LastName     *string            `json:"lastName" validate:"required"`
// 	Password     *string            `json:"password" validate:"required,min=6,max=100"`
// 	Email        *string            `json:"email" validate:"email,required"`
// 	PhoneNumber  *string            `json:"phoneNumber" validate:"required"`
// 	Avator       *string            `json:"avator"`
// 	AccessToken  *string            `json:"accessToken"`
// 	RefreshToken *string            `json:"refreshToken"`
// 	UserType     *string            `json:"userType" validate:"required,eq=ADMIN|eq=USER"`
// 	CreatedAt    time.Time          `json:"createdAt"`
// 	UpdatedAt    time.Time          `json:"updatedAt"`
// 	UserId       string             `json:"userId"`
// }

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastName"`
	Password     *string            `json:"password" validate:"required,min=6,max=100"`
	Email        *string            `json:"email" validate:"email,required"`
	PhoneNumber  *string            `json:"phoneNumber"`
	Avator       *string            `json:"avator"`
	AccessToken  *string            `json:"accessToken"`
	RefreshToken *string            `json:"refreshToken"`
	UserType     *string            `json:"userType"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserId       string             `json:"userId"`
}
