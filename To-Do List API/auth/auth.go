package auth

import (
	"time"

	"github.cim/RohithBN/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(u *types.User) (string, error) {

	claims := jwt.MapClaims{
		"userId":   u.UserID,
		"username": u.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString([]byte("secret_key_placeholder"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func VerifyToken( tokenString string) (*jwt.Token, error) {
	token,err:= jwt.Parse(tokenString, func(token *jwt.Token) (any , error){
		return []byte("secret_key_placeholder"), nil 
	})
	
	if err!= nil{
		return nil, err
	}
	if token.Valid{
		return token, nil
	} else {
		return nil, nil
	}
}
