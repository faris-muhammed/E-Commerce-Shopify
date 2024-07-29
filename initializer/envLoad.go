package initializer

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Envload() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
}
