package main

import (
	"fmt"
	"go-vieon-content-crawling-api/cli/configs"
	"os"
)

func main() {
	if app, err := configs.NewUrfaveCliClient().Instance(); err != nil {
		panic("Error while initialize cli app: " + err.Error())
	} else if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: " + err.Error())
	}
}
