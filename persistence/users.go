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

func (ur *UserRepository) Create(user domain.User) (domain.User, error) {
	result := ur.db.Create(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (ur *UserRepository) Get(username string) (domain.User, error) {
	var user domain.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, errors.New(fmt.Sprintf("User %v not found", username))
	}
	if err != nil {
		return domain.User{}, errors.New("DB error")
	}
	return user, nil
}

func (ur *UserRepository) Delete(user domain.User) error {
	result := ur.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
