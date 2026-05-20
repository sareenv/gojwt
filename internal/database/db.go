package database

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// DBConfig / A struct to hold database configuration parameters.
type DBConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"pguser"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB" envDefault:"pgappdb"`
}

func (c *DBConfig) String() string {
	return "DBConfig{Host: " + c.Host + ", Port: " + fmt.Sprintf("%d", c.Port) + ", User: " + c.User + ", DBName: " + c.DBName + "}"
}

func (c *DBConfig) ConnectionURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// LoadDBConfig / A type alias for a string representing the path to the environment file.
func LoadDBConfig() (*DBConfig, error) {
	// Loads the environment variables from the specified path and returns a DBConfig instance.
	var cfg DBConfig
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
