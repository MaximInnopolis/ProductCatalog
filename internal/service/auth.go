package service

import (
	"context"
	"github.com/MaximInnopolis/ProductCatalog/internal/logger"
	"github.com/MaximInnopolis/ProductCatalog/internal/model"
	"github.com/MaximInnopolis/ProductCatalog/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type AuthService struct {
	repo repository.Authorization
}

func (s *AuthService) CreateUser(ctx context.Context, user *model.User) error {
	var err error
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateToken(ctx context.Context, user *model.User) (string, error) {
	dbUser, err := s.repo.GetUser(ctx, user)
	if err != nil {
		return "", err
	}
	// Generate JWT token
	token, err := generateJWT(dbUser)
	if err != nil {
		logger.Printf(ctx, "Error generating JWT token: %s", err)
		return "", err
	}

	logger.Printf(ctx, "Successfully logged in")

	return token, nil
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// generateJWT generates JWT token for user with additional claims
func generateJWT(user *model.User) (string, error) {
	// Create new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set standard claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["sub"] = user.Username

	// Add additional claims
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generatePasswordHash(password string) (string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %s", err)
		return "", err
	}
	return string(hashedPassword), nil
}
