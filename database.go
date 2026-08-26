package hazart

import (
	"fmt"
	"strings"
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
