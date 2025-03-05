// backend/internal/repository/user_repository.go

package repository

import (
	"gosveltekit/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByResetToken(token string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("reset_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByRefreshToken(token string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("refresh_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUsersWithResetTokens finds all users that currently have reset tokens
func (r *UserRepository) FindUsersWithResetTokens() ([]*models.User, error) {
	var users []*models.User
	// Find all users where reset_token is not empty
	if err := r.db.Where("reset_token <> ''").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	// Validate required fields
	if user.Email == "" {
		return gorm.ErrInvalidField
	}
	if user.Username == "" {
		return gorm.ErrInvalidField
	}
	if user.DisplayName == "" {
		return gorm.ErrInvalidField
	}
	return r.db.Save(user).Error
}
