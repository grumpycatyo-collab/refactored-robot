package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func CreateToken(pass string, hmacSampleSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": pass,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string, hmacSampleSecret string) error {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret for validation
		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		return err
	}

	// Check if the token is valid
	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

/*
const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	token, err := utils.CreateToken(user.Id, hmacSampleSecret) // You need to define this function
	if err != nil {
		return "", err
	}
*/
