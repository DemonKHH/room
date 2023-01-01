package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("userType")
	uid := c.GetString("userId")
	err = nil
	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resources")
		return err
	}
	return err
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("userType")
	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resources")
	}
	return err
}
