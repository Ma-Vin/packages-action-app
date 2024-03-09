package main

import (
	"log"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service"
)

// Main funtion to execute the actions process
func main() {
	log.Println("Start packages action")

	var loadedConfig = config.ReadConfiguration()

	if loadedConfig == nil {
		log.Fatalln("Packages action failed")
	}

	packages, err := service.GetAndPrintUserPackages(loadedConfig)
	if err != nil {
		log.Fatalln(err)
	}
	if len(*packages) > 1 {
		_, err = service.GetAndPrintUserPackage((*packages)[0].Name, loadedConfig)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("Packages action done")
}
