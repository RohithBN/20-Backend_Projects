package middleware

import (
	"strings"

	"github.cim/RohithBN/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
		token,err:= auth.VerifyToken(tokenString);
		if err!= nil{
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		
		if token == nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user",token)
		c.Next()

	}
}
