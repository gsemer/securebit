package domain

import "github.com/google/uuid"

type UserPayload struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username       string    `gorm:"type:text;not null;unique" json:"username"`
	HashedPassword string    `gorm:"type:text;not null" json:"-"`
}

type UserRepository interface {
	Create(authUser AuthUser) error
}
