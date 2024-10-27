package database

import (
	"fmt"

	"github.com/malytinKonstantin/go-messenger-mono/auth-service/internal/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		viper.GetString("DATABASE_NAME"),
		viper.GetString("DATABASE_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseDatabaseConnection(db *gorm.DB) error {
	postgresDB, err := db.DB()
	if err != nil {
		return err
	}
	return postgresDB.Close()
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.UserCredentials{}, &models.OauthAccount{}, &models.ResetPasswordToken{})
}
