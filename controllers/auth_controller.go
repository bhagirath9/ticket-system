package controllers

import (
	"errors"
	"net/http"
	"ticket-system/models"
	"ticket-system/services"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication HTTP requests.
type AuthController struct {
	authService services.AuthService
}

// NewAuthController returns a new instance of AuthController.
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register processes requests for new user registration.
// POST /auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	// Validates fields (required, email format, password min length)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.authService.Register(&req)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User Registered Successfully"})
}

// Login authenticates a user and returns a JWT token.
// POST /auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	// Validates required email (email format check) and password
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.authService.Login(&req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}
