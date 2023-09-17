package service

import (
	"time"

	"github.com/SawitProRecruitment/UserService/repository"
)

// AuthService represents the authentication service interface.
// Intended as a provider
type AuthService interface {
	// CreateToken generates a token for the given user and returns it.
	// It takes a user object as input and returns the generated token as a string
	// or an error if token generation fails.
	CreateToken(user *repository.User, now time.Time) (string, error)

	// ValidateToken validates the provided token and returns the user's ID.
	// It takes a token string as input and returns the user's ID as an int64 if
	// the token is valid, or an error if validation fails.
	ValidateToken(token string) (int64, error)
}
