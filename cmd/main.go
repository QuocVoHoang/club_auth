package main

import (
	"fmt"
	"os"

	"github.com/your-org/go-base/internal/app"
	"github.com/your-org/go-base/internal/config"
)

func main() {
	// Load config
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(cfg); err != nil {
		fmt.Printf("Failed to run application: %v\n", err)
		os.Exit(1)
	}
}
