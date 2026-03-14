package main

import (
	"errors"
	"fmt"
	"os"

	"gosveltekit/internal/bootstrap"
	"gosveltekit/internal/config"
	"gosveltekit/internal/models"
	"gosveltekit/internal/validation"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	identifier := os.Getenv("ADMIN_IDENTIFIER")
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	displayName := os.Getenv("ADMIN_DISPLAY_NAME")
	if displayName == "" {
		displayName = "Administrator"
	}

	if err := validateInput(identifier, email, password, displayName); err != nil {
		exitf("invalid admin seed input: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		exitf("failed to load config: %v", err)
	}

	db, err := bootstrap.OpenGorm(cfg)
	if err != nil {
		exitf("failed to open database: %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		exitf("failed to hash password: %v", err)
	}

	user, created, err := upsertAdmin(db, identifier, email, displayName, string(passwordHash))
	if err != nil {
		exitf("failed to seed admin: %v", err)
	}

	action := "updated"
	if created {
		action = "created"
	}

	fmt.Printf("%s admin user %s (%s)\n", action, user.Username, user.Email)
}

func validateInput(identifier, email, password, displayName string) error {
	if err := validation.ValidateUsername(identifier); err != nil {
		return err
	}
	if err := validation.ValidateEmail(email); err != nil {
		return err
	}
	if err := validation.ValidatePassword(password, identifier); err != nil {
		return err
	}
	if err := validation.ValidateDisplayName(displayName); err != nil {
		return err
	}

	return nil
}

func upsertAdmin(
	db *gorm.DB,
	identifier, email, displayName, passwordHash string,
) (*models.User, bool, error) {
	var user models.User
	created := false

	err := db.Where("username = ?", identifier).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Where("email = ?", email).First(&user).Error
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, err
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{
			Username:      identifier,
			Email:         email,
			DisplayName:   displayName,
			PasswordHash:  passwordHash,
			Role:          "admin",
			Active:        true,
			EmailVerified: true,
		}

		if err := db.Create(&user).Error; err != nil {
			return nil, false, err
		}

		created = true
		return &user, created, nil
	}

	user.Username = identifier
	user.Email = email
	user.DisplayName = displayName
	user.PasswordHash = passwordHash
	user.Role = "admin"
	user.Active = true
	user.EmailVerified = true

	if err := db.Save(&user).Error; err != nil {
		return nil, false, err
	}

	return &user, created, nil
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
