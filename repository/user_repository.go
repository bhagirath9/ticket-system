package repository

import (
	"ticket-system/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for database operations related to Users.
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user record into the database.
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail searches for a user record matching the given email.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID retrieves a user record by their primary key ID.
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
