package main

import (
	"api"
	"fmt"
)

func main() {
	fmt.Println("Fader starting...")

	fmt.Println("Init config...")
	initConfig()

	fmt.Println("Init elasticsearch...")
	initElasticSearch()

	fmt.Println("Init stores...")
	initStroes()

	fmt.Println("Api...")
	api.Run()
}
