package persistence

import (
	"securebit/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(authUser domain.AuthUser) error {
	result := ur.db.Create(&authUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
