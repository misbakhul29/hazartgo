package hazart

import (
	"testing"
)

func TestDBConfigBuildDSN(t *testing.T) {
	tests := []struct {
		config   DBConfig
		expected string
	}{
		{
			config: DBConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "secretpassword",
				DBName:   "mydb",
				SSLMode:  "disable",
			},
			expected: "host=localhost user=postgres password=secretpassword dbname=mydb port=5432 sslmode=disable",
		},
		{
			config: DBConfig{
				Driver:   "mysql",
				Host:     "127.0.0.1",
				Port:     3306,
				User:     "root",
				Password: "rootpassword",
				DBName:   "mysqldb",
			},
			expected: "root:rootpassword@tcp(127.0.0.1:3306)/mysqldb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			config: DBConfig{
				Driver: "sqlite",
				File:   "data.db",
			},
			expected: "data.db",
		},
	}

	for _, tt := range tests {
		dsn := tt.config.BuildDSN()
		if dsn != tt.expected {
			t.Errorf("BuildDSN() = %q, expected %q", dsn, tt.expected)
		}
	}
}
