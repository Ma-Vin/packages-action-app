package main

import (
	"log"

	"github.com/ma-vin/packages-action/config"
)

// Main funtion to execute the actions process
func main() {
	log.Println("Start packages action")

	_, err := config.ReadConfiguration()

	checkkError(err)

	log.Println("Packages action done")
}

func checkkError(err error) {
	if err != nil {
		log.Fatalf("Packages action failed: %s", err)
	}
}
