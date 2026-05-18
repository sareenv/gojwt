package main

import (
	"fmt"

	"github.com/sareenv/gojwt/internal/database"
)

func main() {
	// Load database configuration.
	cfg, err := database.LoadDBConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database configuration loaded:", cfg)
}
