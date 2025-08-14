package domain

import "github.com/google/uuid"

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username       string    `gorm:"type:text;not null;unique" json:"username"`
	HashedPassword string    `gorm:"type:text;not null" json:"-"`
}

type UserRepository interface {
	Create(user User) (User, error)
	Get(username string) (User, error)
	Delete(use User) error
}
