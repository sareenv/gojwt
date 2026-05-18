package database

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// / A struct to hold database configuration parameters.
type DBConfig struct {
	Host     string `env:"POSTGRES_USER" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"pguser"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB" envDefault:"pgappdb"`
}

func (c *DBConfig) String() string {
	return "DBConfig{Host: " + c.Host + ", Port: " + fmt.Sprintf("%d", c.Port) + ", User: " + c.User + ", DBName: " + c.DBName + "}"
}

// / NewDBConfig creates a new instance of DBConfig with the provided parameters.
func NewDBConfig(host string, port int, user, password, dbName string) *DBConfig {
	return &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}
}

// / A type alias for a string representing the path to the environment file.
func LoadDBConfig() (*DBConfig, error) {
	// Loads the enviornment variables from the specified path and returns a DBConfig instance.
	var cfg DBConfig
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
