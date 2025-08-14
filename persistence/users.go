package persistence

import (
	"errors"
	"fmt"
	"securebit/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(authUser domain.AuthUser) (domain.AuthUser, error) {
	result := ur.db.Create(&authUser)
	if result.Error != nil {
		return domain.AuthUser{}, result.Error
	}
	return authUser, nil
}

func (ur *UserRepository) Get(username string) (domain.AuthUser, error) {
	var user domain.AuthUser
	err := ur.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.AuthUser{}, errors.New(fmt.Sprintf("User %v not found", username))
	}
	if err != nil {
		return domain.AuthUser{}, errors.New("DB error")
	}
	return user, nil
}

func (ur *UserRepository) Delete(authUser domain.AuthUser) error {
	result := ur.db.Delete(&authUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
