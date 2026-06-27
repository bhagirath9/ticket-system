package services

import (
	"errors"
	"ticket-system/config"
	"ticket-system/models"
	"ticket-system/repository"
	"ticket-system/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Sentinel errors for auth flows.
var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// AuthService defines the interface for auth-related business logic.
type AuthService interface {
	Register(req *models.RegisterRequest) error
	Login(req *models.LoginRequest) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

// NewAuthService returns a new instance of AuthService.
func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// Register registers a new user by checking duplicate email, hashing the password, and storing.
func (s *authService) Register(req *models.RegisterRequest) error {
	// Check if a user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return ErrEmailExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

// Login validates user credentials and returns a signed JWT token on success.
func (s *authService) Login(req *models.LoginRequest) (string, error) {
	// Look up the user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	// Validate the password against the stored bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate and sign the JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
