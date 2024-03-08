package models

import (
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// IsTokenValid checks if token is valid
func IsTokenValid(db *sql.DB, tokenString string) (bool, error) {

	// Retrieve token from database
	var dbToken string
	query := "SELECT token FROM users WHERE token = ?"
	err := db.QueryRow(query, tokenString).Scan(&dbToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Println("Error retrieving token from database:", err)
			return false, nil
		}
		return false, err
	}

	// Check token validity
	validToken, err := checkToken(tokenString)
	if err != nil {
		logger.Println("Error checking token validity:", err)
		return false, err
	}
	if !validToken {
		return false, errors.New("invalid token")
	}

	logger.Println("Token is valid")
	return true, nil
}

// RegisterUser registers new user in database
func RegisterUser(db *sql.DB, user *User) error {
	// Check if username already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		logger.Println("Error checking if username exists:", err)
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Println("Error hashing password:", err)
		return err
	}

	// Insert new user in database
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, hashedPassword)
	if err != nil {
		logger.Println("Error inserting new user in database:", err)
		return err
	}

	logger.Println("Successfully registered")

	return nil
}

// LoginUser login user and generate JWT token
func LoginUser(db *sql.DB, user *User) (string, error) {
	var dbUser User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", user.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		logger.Println("Error retrieving user from database:", err)
		return "", err
	}

	// Compare stored password hash with provided password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token, err := generateJWT(&dbUser)
	if err != nil {
		logger.Println("Error generating JWT token:", err)
		return "", err
	}

	logger.Println("Successfully logged in")

	return token, nil
}

func checkToken(tokenString string) (bool, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return false, err
	}

	// Check if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Println("Token is invalid")
		return false, nil
	}

	// Check if expiration claim exists and validate it
	expiration, ok := claims["exp"].(float64)
	if !ok {
		logger.Println("no expiration claim found")
		return false, nil
	}

	if int64(expiration) < time.Now().Unix() {
		logger.Println("Token has expired")
		return false, nil
	}

	return true, nil
}

// GenerateJWT generates JWT token for user with additional claims
func generateJWT(user *User) (string, error) {
	// Get secret key from environment variable
	secretKey := os.Getenv("SECRET_KEY")

	// Create new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set standard claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["sub"] = user.Username

	// Add additional claims
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
