package auth

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthService provides authentication services
type AuthService struct {
	DB       *gorm.DB
	FireAuth *auth.Client
}

// Login authenticates a user with the provided credentials and returns a Firebase custom token
func (s *AuthService) Login(email, password string) (string, error) {
	// Get the user from the database
	var user User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("user with email does not exist")
		}
		log.Printf("failed to get user by email from database: %v", err)
		return "", errors.New("internal server error")
	}

	// Check if the provided password matches the user's password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("incorrect password")
	}

	// Generate a Firebase custom token for the user
	token, err := s.FireAuth.CustomToken(context.Background(), user.ID)
	if err != nil {
		log.Printf("failed to generate custom token: %v", err)
		return "", errors.New("internal server error")
	}

	return token, nil
}

// Register creates a new user with the provided credentials and returns token
func (s *AuthService) Register(email, password string) (string, error) {
	// Check if the user with the email already exists
	var user User
	if err := s.DB.Raw("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("failed to get user by email from database: %v", err)
			return "", errors.New("internal server error")
		}
	}

	if user.ID != "" {
		return "", errors.New("user with email already exists")
	}

	// Generate a UUID for the new user
	uid := uuid.New().String()

	// Generate a hash of the user's password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return "", errors.New("internal server error")
	}

	// Create a new user in the database
	if err := s.DB.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", uid, email, hashedPassword).Error; err != nil {
		log.Printf("failed to insert user into database: %v", err)
		return "", errors.New("internal server error")
	}

	// Create a custom token for the user using the Firebase Admin SDK
	customToken, err := s.FireAuth.CustomToken(context.Background(), uid)
	if err != nil {
		log.Printf("failed to create custom token for user: %v", err)
		return "", errors.New("internal server error")
	}

	return customToken, nil
}
