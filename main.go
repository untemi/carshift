package main

import (
	"fmt"
	"os"

	"github.com/untemi/carshift/cmd"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No action provided")
		fmt.Println("Actions : serve, setup")
		return
	}

	switch os.Args[1] {
	case "serve":
		cmd.Serve()
	case "setup":
		cmd.Setup()
	default:
		fmt.Printf("Unknown action : %s\n", os.Args[1])
		fmt.Println("Actions : serve, setup")
	}
}
