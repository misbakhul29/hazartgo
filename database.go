package hazart

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig holds database connection parameters and auto-migration models
type DBConfig struct {
	Driver      string // "postgres", "mysql", "sqlite"
	Host        string // Default "localhost"
	Port        int    // Default 5432 (postgres) or 3306 (mysql)
	User        string
	Password    string
	DBName      string
	SSLMode     string // Default "disable" for postgres
	File        string // Default "app.db" for sqlite
	AutoMigrate []any  // Models to auto migrate e.g. []any{&Product{}, &User{}}
}

// BuildDSN constructs the database DSN connection string automatically based on config
func (c *DBConfig) BuildDSN() string {
	driver := strings.ToLower(c.Driver)

	switch driver {
	case "mysql":
		port := c.Port
		if port == 0 {
			port = 3306
		}
		host := c.Host
		if host == "" {
			host = "127.0.0.1"
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, host, port, c.DBName)

	case "sqlite", "sqlite3":
		file := c.File
		if file == "" {
			file = "app.db"
		}
		return file

	default: // postgres
		port := c.Port
		if port == 0 {
			port = 5432
		}
		host := c.Host
		if host == "" {
			host = "localhost"
		}
		sslmode := c.SSLMode
		if sslmode == "" {
			sslmode = "disable"
		}
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			host, c.User, c.Password, c.DBName, port, sslmode)
	}
}

// OpenDB opens a database connection using GORM based on the provided DBConfig.
// Driver supports "sqlite", "postgres", and "mysql".
// If AutoMigrate models are defined in config, it automatically migrates them.
func OpenDB(config DBConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(config.Driver) {
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(config.BuildDSN())
	case "postgres", "postgresql":
		dialector = postgres.Open(config.BuildDSN())
	case "mysql":
		dialector = mysql.Open(config.BuildDSN())
	default:
		// Default fallback to sqlite if driver is unspecified or sqlite
		dialector = sqlite.Open(config.BuildDSN())
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", config.Driver, err)
	}

	if len(config.AutoMigrate) > 0 {
		if err := AutoMigrate(db, config.AutoMigrate...); err != nil {
			return db, fmt.Errorf("database connected but auto-migration failed: %w", err)
		}
	}

	return db, nil
}

// DatabaseConnection is an alias for OpenDB for developer ergonomics.
func DatabaseConnection(config DBConfig) (*gorm.DB, error) {
	return OpenDB(config)
}

// AutoMigrate migrates schema models for the given GORM DB instance.
func AutoMigrate(db *gorm.DB, models ...any) error {
	if db == nil {
		return fmt.Errorf("gorm db instance is nil")
	}
	if len(models) == 0 {
		return nil
	}
	return db.AutoMigrate(models...)
}

