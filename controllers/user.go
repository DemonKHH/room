package controllers

import (
	"net/http"

	"wmt/helpers"
	response "wmt/pkg/common/response"
	db "wmt/service/mongo"
	serviceUser "wmt/service/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoClient = db.GetMongoClient()
var userCollection *mongo.Collection = db.OpenCollection(mongoClient, "users")
var validate = validator.New()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
			return
		}
		user, err := serviceUser.GetUser(userId)
		if err != nil {
			// c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
		}
		// c.JSON(http.StatusOK, user)
		c.JSON(http.StatusOK, response.ResponseMsg{
			Code: 0,
			Msg:  "登录成功",
			Data: user,
		})
	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var users []bson.M
		if err = helpers.CheckUserType(c, "ADMIN"); err != nil {
			// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
			return
		}
		users, err = serviceUser.GetUsers()
		if err != nil {
			// c.JSON(http.StatusOK, err.Error())
			c.JSON(http.StatusOK, response.FailMsg(
				err.Error(),
			))
			return
		}
		// c.JSON(http.StatusOK, users)
		c.JSON(http.StatusOK, response.ResponseMsg{
			Code: 0,
			Msg:  "登录成功",
			Data: users,
		})
	}
}
