package model

import (
	"context"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// IsTokenValid checks if token is valid
func IsTokenValid(ctx context.Context, tokenString string) (bool, error) {

	// Check token validity
	validToken, err := CheckToken(ctx, tokenString)
	if err != nil || !validToken {
		logger.Printf(ctx, "Error checking token validity: %s", err)
		return false, errors.New("invalid token")
	}

	logger.Printf(ctx, "Token is valid")
	return true, nil
}

func CheckToken(ctx context.Context, tokenString string) (bool, error) {

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return false, err
	}

	// Check if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Printf(ctx, "Token is invalid")
		return false, nil
	}

	// Check if expiration claim exists and validate it
	expiration, ok := claims["exp"].(float64)
	if !ok {
		logger.Printf(ctx, "no expiration claim found")
		return false, nil
	}

	if int64(expiration) < time.Now().Unix() {
		logger.Printf(ctx, "Token has expired")
		return false, nil
	}

	return true, nil
}
