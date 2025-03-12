package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/untemi/carshift/internal"
	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/handler"
	"github.com/untemi/carshift/internal/misc"
)

func Serve() {
	if !misc.IsFileExists(".up") {
		fmt.Println("Server not set up, run 'carshift setup'")
		fmt.Println("(if you think this is wrong run 'touch .up')")
		return
	}

	var serverConfig internal.ServerConfig
	ctx := context.Background()

	// Init
	flagParse(&serverConfig)

	// Database
	closeDB, err := db.Init(ctx)
	if err != nil {
		log.Fatalf("DB: Error initialising database, %v", err)
	}
	defer closeDB()

	// Handlers
	closeHDB, err := handler.Init()
	if err != nil {
		log.Fatalf("SM: Error initialising session manager, %v", err)
	}
	defer closeHDB()

	// Starting Server
	internal.Start(serverConfig)
}
