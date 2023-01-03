package config

import (
	"chigitaction/models"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", viper.GetString("DB_HOST"), viper.GetString("DB_PORT"), viper.GetString("DB_USER"), viper.GetString("DB_NAME"), viper.GetString("DB_PASS"))

	db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDb, err := db.DB()

	err = sqlDb.Ping()

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Activity{}, &models.Todo{})

	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(10)

	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return db, nil
}
