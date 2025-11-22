package main

import (
	"log"

	"github.com/Vitaliy-X/avito-internship-autumn-2025/internal/app"
)

func main() {
	application := app.New()

	if err := application.Run(); err != nil {
		log.Fatalf("Application stopped: %v", err)
	}
}
