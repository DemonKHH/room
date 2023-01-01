package helpers

import (
	"context"

	"log"
	"os"
	"time"

	db "wmt/service/mongo"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	// LastName  string
	Uid string
	// UserType  string
	TokenType string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = db.OpenCollection(db.GetMongoClient(), "users")

func GenerateAllToken(email string, firstName string, uid string) (signedAccessToken string, signedRefreshToken string, err error) {
	accessClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		// LastName:  lastName,
		// UserType:  userType,
		Uid:       uid,
		TokenType: "accessToken",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(2000)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Uid:       uid,
		TokenType: "refreshToken",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*7)).Unix(),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Panic(err)
		return
	}

	return accessToken, refreshToken, err

}

func UpdateAllTokens(signedAccessToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	// var updateObj primitive.D

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	id, _ := primitive.ObjectIDFromHex(userId)
	_, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "accesstoken", Value: signedAccessToken}, {Key: "refreshtoken", Value: signedRefreshToken}, {Key: "updatedat", Value: updatedAt}}},
		},
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = err.Error()
		return
	}

	return claims, msg
}
