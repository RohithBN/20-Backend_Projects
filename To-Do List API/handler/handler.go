package handler

import (
	"fmt"

	"github.cim/RohithBN/auth"
	"github.cim/RohithBN/types"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	var newUser *types.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if newUser.Username == "" || newUser.Password == "" {
		c.JSON(400, gin.H{"error": "Username and Password are required"})
		return
	}

	// TODO: store user in database

	c.JSON(201, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var loginUser *types.User

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if loginUser.Username == "" || loginUser.Password == "" {
		c.JSON(400, gin.H{"error": "Username and Password are required"})
		return
	}

	// TODO: verify user in database

	verifiedUser := &types.User{
		UserID:   1,
		Username: loginUser.Username,
		Password: loginUser.Password,
	}

	fmt.Println("User", verifiedUser)

	// generate a jwt token

	tokenString,err:= auth.GenerateToken(verifiedUser)
	if err!= nil{
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	fmt.Println("TOKEN", tokenString)

	c.JSON(200, gin.H{"message": "Login successful", "token": tokenString})
}
