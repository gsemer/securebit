package persistence

import (
	"fmt"
	"securebit/config"
	"securebit/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Username string
	Password string
	DBName   string
	Host     string
	Port     string
}

func LoadPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Username: config.GetEnv("POSTGRES_USERNAME", ""),
		Password: config.GetEnv("POSTGRES_PASSWORD", ""),
		DBName:   config.GetEnv("POSTGRES_DB", ""),
		Host:     config.GetEnv("POSTGRES_IP", ""),
		Port:     config.GetEnv("POSTGRES_PORT", ""),
	}
}

func (config *PostgresConfig) Init() (*gorm.DB, error) {
	source := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := gorm.Open(postgres.Open(source), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&domain.AuthUser{})

	return db, nil
}
