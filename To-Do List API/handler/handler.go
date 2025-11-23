package handler

import (
	"fmt"

	"github.cim/RohithBN/auth"
	"github.cim/RohithBN/lib"
	"github.cim/RohithBN/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword , err:= bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err!= nil{
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = lib.DB.Exec(query, newUser.Username, string(hashedPassword))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}
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
	query:= "SELECT user_id, password FROM users WHERE username = ?"
	row:= lib.DB.QueryRow(query, loginUser.Username)

	var storedHashedPassword string
	var userID int
	
	err:= row.Scan(&userID, &storedHashedPassword)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(loginUser.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	verifiedUser := &types.User{
		UserID:   uint(userID),
		Username: loginUser.Username,
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
