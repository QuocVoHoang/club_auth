package database

import (
	"context"
	"fmt"

	"github.com/your-org/go-base/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SSLMode      string
	ENV          string
}

func InitDatabase(c Config) (*gorm.DB, error) {
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Host,
		c.User,
		c.Password,
		c.DatabaseName,
		c.Port,
		c.SSLMode,
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if c.ENV != config.ENVProduction {
		conn = conn.Debug()
	}

	return conn, nil
}

func Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}
