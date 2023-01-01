package middleware

import (
	"net/http"
	"strings"

	"wmt/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			clientToken = c.Request.Header.Get("token")
		} else if strings.HasPrefix(clientToken, "Bearer ") {
			reqToken := c.Request.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, "Bearer ")
			clientToken = splitToken[1]
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid authorization token"})
			c.Abort()
			return
		}

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no Authorization header provided"})
			c.Abort()
			return
		}
		// handle access token
		claims, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		if claims.TokenType == "accessToken" {
			c.Set("email", claims.Email)
			c.Set("firstName", claims.FirstName)
			// c.Set("lastName", claims.LastName)
			c.Set("userId", claims.Uid)
			// c.Set("userType", claims.UserType)
			c.Next()
		} else if claims.TokenType == "refreshToken" {
			c.Set("tokenType", claims.TokenType)
			c.Set("userId", claims.Uid)
			c.Next()
		}

	}
}
