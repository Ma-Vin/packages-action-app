package main

import (
	"log"

	"github.com/ma-vin/packages-action/config"
)

// Main funtion to execute the actions process
func main() {
	log.Println("Start packages action")

	var loadedConfig = config.ReadConfiguration()

	if loadedConfig == nil {
		log.Fatalln("Packages action failed")
	}

	log.Println("Packages action done")
}

func checkkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
