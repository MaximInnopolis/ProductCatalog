package models_test

import (
	"context"
	"database/sql"
	"github.com/MaximInnopolis/ProductCatalog/internal/models"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestIsTokenValid(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "auth/register")

	// Create temporary database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create users table
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		t.Fatalf("Error creating users table: %v", err)
	}

	// Create user
	user := &models.User{Username: "testuser", Password: "testpassword"}
	if err := models.RegisterUser(ctx, db, user); err != nil {
		t.Fatalf("Error registering user: %v", err)
	}

	// Generate JWT token
	token, err := models.LoginUser(ctx, db, user)
	if err != nil {
		t.Fatalf("Error generating JWT token: %v", err)
	}

	// Check token validity
	valid, err := models.IsTokenValid(ctx, token)
	if err != nil {
		t.Fatalf("Error checking token validity: %v", err)
	}

	if !valid {
		t.Error("Expected  valid token, but got invalid one")
	}
}

func TestRegisterUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "auth/register")
	// Create temporary database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create users table
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		t.Fatalf("Error creating users table: %v", err)
	}

	// Create user
	user := &models.User{Username: "testuser", Password: "testpassword"}

	// Register user
	err = models.RegisterUser(ctx, db, user)
	if err != nil {
		t.Fatalf("Error registering user: %v", err)
	}

	// Check if registered user exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		t.Fatalf("Error executing query: %v", err)
	}

	if count != 1 {
		t.Error("Expected user to be registered in database, but no records found")
	}
}

func TestLoginUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "auth/login")

	// Create temporary database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create users table
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT)")
	if err != nil {
		t.Fatalf("Error creating users table: %v", err)
	}

	// Create user
	user := &models.User{Username: "testuser", Password: "testpassword"}

	// Register user
	err = models.RegisterUser(ctx, db, user)
	if err != nil {
		t.Fatalf("Error registering user: %v", err)
	}

	// Attempt to login user
	_, err = models.LoginUser(ctx, db, user)
	if err != nil {
		t.Fatalf("Error attempting to login user: %v", err)
	}
}

func TestCheckTokenValid(t *testing.T) {
	ctx := context.WithValue(context.Background(), "endpoint", "auth/login")
	// Generate valid token
	user := &models.User{ID: 123, Username: "testuser"}
	tokenString, err := models.GenerateJWT(user)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Check validity of token
	valid, err := models.CheckToken(ctx, tokenString)
	if err != nil {
		t.Fatalf("Error checking token: %v", err)
	}
	if !valid {
		t.Error("Expected token to be valid, got invalid")
	}
}

func TestGenerateJWT(t *testing.T) {
	// Create user
	user := &models.User{ID: 123, Username: "testuser"}

	// Generate JWT token
	tokenString, err := models.GenerateJWT(user)
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Check if token string is not empty
	if tokenString == "" {
		t.Error("Expected non-empty token string, got empty")
	}
}
