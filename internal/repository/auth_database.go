package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthDatabase struct {
	db *sql.DB
}

// CreateUser create new user in database
func (r *AuthDatabase) CreateUser(ctx context.Context, user *model.User) error {
	// Check if username already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		logger.Printf(ctx, "Error checking if username exists: %s", err)
		return err
	}
	if count > 0 {
		logger.Printf(ctx, "Username already exists")
		return errors.New("username already exists")
	}

	// Insert new user in database
	_, err = r.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		logger.Printf(ctx, "Error inserting new user in database: %s", err)
		return err
	}

	logger.Printf(ctx, "Successfully registered")

	return nil
}

// GetUser returns user
func (r *AuthDatabase) GetUser(ctx context.Context, user *model.User) (*model.User, error) {
	var dbUser model.User
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", user.Username).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		logger.Printf(ctx, "Error retrieving user from database: %s", errors.New("not registered yet"))
		return nil, err
	}

	// Compare stored password hash with provided password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &dbUser, nil
}

func NewAuthDatabase(db *sql.DB) *AuthDatabase {
	return &AuthDatabase{db: db}
}
