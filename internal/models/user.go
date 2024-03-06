package models

import (
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// AuthenticateUser authenticates user
func AuthenticateUser(db *sql.DB, user *User) (bool, error) {

	// Retrieve user data from the database
	var dbUser User
	err := db.QueryRow("SELECT id, username, password, token FROM users WHERE username = ?", user.Username).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password, &dbUser.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Println("Error retrieving user data from database:", err)
			return false, nil
		}
		return false, err
	}

	// Check if retrieved token is NULL
	if dbUser.Token == "" {
		return false, errors.New("token not found")
	}

	// Check token validity
	validToken, err := checkToken(user.Token, dbUser.Token)
	if err != nil {
		logger.Println("Error checking token validity:", err)
		return false, err
	}
	if !validToken {
		return false, errors.New("invalid token")
	}

	// Retrieve hashed password from the database
	var dbPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&dbPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.Println("Error retrieving hashed password from database:", err)
		return false, err
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(user.Password))
	if err != nil {
		return false, nil
	}

	logger.Println("Authentication successful")
	return true, nil
}

// RegisterUser registers new user in database
func RegisterUser(db *sql.DB, user *User) (string, error) {
	// Check if the username already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		logger.Println("Error checking if username exists:", err)
		return "", err
	}
	if count > 0 {
		return "", errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Println("Error hashing password:", err)
		return "", err
	}

	// Generate JWT token
	token, err := generateJWT(user.Password)
	if err != nil {
		logger.Println("Error generating JWT token:", err)
		return "", err
	}

	// Insert new user in database
	_, err = db.Exec("INSERT INTO users (username, password, token) VALUES (?, ?, ?)", user.Username, hashedPassword, token)
	if err != nil {
		logger.Println("Error inserting new user in database:", err)
		return "", err
	}

	logger.Println("Successfully registered")

	return token, nil
}

func checkToken(tokenString string, dbToken string) (bool, error) {
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
	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return false, nil
	}

	// Check if tokens match
	if tokenString != dbToken {
		return false, nil
	}

	return true, nil
}

// GenerateJWT generates JWT token for user
func generateJWT(password string) (string, error) {
	// Token creation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password": password,
	})

	// Signing token with secret key
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
