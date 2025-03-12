package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/handler"
	"github.com/untemi/carshift/internal/misc"
)

func Setup() {
	// Check if already done
	if misc.IsFileExists(".up") {
		fmt.Println("Already done.")
		fmt.Println("(if you think this is wrong remove '.up' file.)")
		return
	}

	ctx := context.Background()

	// Database
	if err := db.Setup(ctx); err != nil {
		log.Fatalf("DB: Error setting up database, %v", err)
	}

	// Handlers
	if err := handler.Setup(ctx); err != nil {
		log.Fatalf("DB: Error setting up database, %v", err)
	}

	os.WriteFile(".up",
		[]byte("DO NOT REMOVE UNLESS YOU KNOW WHAT YOU ARE DOING"),
		0644,
	)
}
