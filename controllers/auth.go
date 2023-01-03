package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"time"

	"wmt/helpers"

	modelUser "wmt/internal/model/user"
	response "wmt/pkg/common/response"
	"wmt/pkg/common/util"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPass string, providedPass string) (passIsValid bool, msg string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(providedPass))
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(userPass, providedPass)
		return false, fmt.Sprint("email or password is incorrect")
	} else {
		return true, ""
	}

}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user modelUser.User
		// var user modelUser.SignupUser

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
			defer cancel()
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusOK, response.FailMsg(
				validationErr.Error(),
			))
			defer cancel()
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusOK, response.FailMsg(
				"error occured while chacking email existance",
			))
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		// count, err = userCollection.CountDocuments(ctx, bson.M{"phonenumber": user.PhoneNumber})
		// defer cancel()
		// if err != nil {
		// 	log.Panic(err)
		// 	c.JSON(http.StatusOK, response.FailMsg(
		// 		"error occured while chacking phone number existance",
		// 	))
		// 	return
		// }

		if count > 0 {
			c.JSON(http.StatusOK, response.FailMsg(
				"This email or phone number already exist",
			))
			return
		}
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		// accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserId)
		accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.FirstName, user.UserId)
		user.AccessToken = &accessToken
		user.RefreshToken = &refreshToken
		avator := util.GetRandomAvator()
		user.Avator = &avator

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprint("User item was not created")
			c.JSON(http.StatusOK, response.FailMsg(
				msg,
			))
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "用户注册成功",
			"data": getLoginResponse(user),
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user modelUser.User
		var foundUser modelUser.User
		// var user modelUser.LoginUser
		// var foundUser modelUser.SignupUser

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
			defer cancel()
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusOK, response.FailMsg(
				"Email or password is incorrect",
			))
			defer cancel()
			return
		}

		passIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		defer cancel()

		if !passIsValid {
			c.JSON(http.StatusOK, response.FailMsg(
				msg,
			))
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusOK, response.FailMsg(
				"user not found",
			))
			return
		}

		// accessToken, refreshToken, _ := helpers.GenerateAllToken(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserId)
		accessToken, refreshToken, _ := helpers.GenerateAllToken(*foundUser.Email, *foundUser.FirstName, foundUser.UserId)
		helpers.UpdateAllTokens(accessToken, refreshToken, foundUser.UserId)

		err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserId}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"code": 0,
		// 	"msg":  "登录成功",
		// 	"data": getLoginResponse(foundUser),
		// })
		c.JSON(http.StatusOK, response.ResponseMsg{
			Code: 0,
			Msg:  "登录成功",
			Data: getLoginResponse(foundUser),
		})
	}
}

type TokenModel struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("tokenType") == "refreshToken" {
			var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			var user modelUser.User
			UserId := c.GetString("UserId")
			id, _ := primitive.ObjectIDFromHex(UserId)
			err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
			if err != nil {
				c.JSON(http.StatusOK, response.FailMsg(
					err.Error(),
				))
				defer cancel()
				return
			}

			// accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserId)
			accessToken, refreshToken, _ := helpers.GenerateAllToken(*user.Email, *user.FirstName, user.UserId)
			helpers.UpdateAllTokens(accessToken, refreshToken, user.UserId)

			var tokenModel = TokenModel{AccessToken: accessToken, RefreshToken: refreshToken}
			c.JSON(http.StatusOK, tokenModel)
			defer cancel()

		} else {
			c.JSON(http.StatusOK, response.FailMsg(
				"invald refresh token",
			))
		}
	}
}

func getLoginResponse(user modelUser.User) response.LoginResponse {
	var loginResponse = response.LoginResponse{}
	b, _ := json.Marshal(&user)
	json.Unmarshal(b, &loginResponse)
	return loginResponse
}
