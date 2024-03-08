package main

import (
	"fmt"
	"os"

	"github.com/ma-vin/packages-action/config"
)

// Main funtion to execute the actions process
func main() {
	fmt.Println("Start packages action")

	var loadedConfig = config.ReadConfiguration()

	if loadedConfig == nil {
		fmt.Println("Packages action failed")
		os.Exit(1)
		return
	}

	fmt.Println("Packages action done")
}
