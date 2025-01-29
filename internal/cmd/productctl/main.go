package main

import (
	"log"

	"github.com/opdev/productctl/internal/cmd/productctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
