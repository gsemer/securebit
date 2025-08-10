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

func (ur *UserRepository) Create(authUser domain.AuthUser) (domain.AuthUser, error) {
	result := ur.db.Create(&authUser)
	if result.Error != nil {
		return domain.AuthUser{}, result.Error
	}
	return authUser, nil
}
