package persistence

import (
	"fmt"
	"securebit/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(username, password, ip, port, name string) (*gorm.DB, error) {
	source := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, ip, port, name)
	db, err := gorm.Open(postgres.Open(source), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&domain.User{})

	return db, nil
}
