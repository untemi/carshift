package cmd

import (
	"flag"

	"github.com/untemi/carshift/internal"
)

func flagParse(c *internal.ServerConfig) {
	flag.StringVar(&c.Address, "a", ":8000", "address")

	flag.Parse()
}
