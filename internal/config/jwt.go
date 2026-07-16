package config

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken creates a JWT token valid for 24 hours
func GenerateToken(userID int, email string) (string, error) {
	token :=  jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email": email,
		"exp":time.Now().Add(24*time.Hour).Unix(),
	})

	return token.SignedString(jwtSecret)
}

// VerifyToken validates a JWT token
func VerifyToken(tokenString string) (int, error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))
	return userID, nil
}
