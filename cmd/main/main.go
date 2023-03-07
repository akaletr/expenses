package main

import (
	"log"

	"cmd/main/main.go/internal/app"
	"cmd/main/main.go/internal/config"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	myApp := app.NewApp(cfg)
	defer func() {
		err = myApp.Stop()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = myApp.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = myApp.Start()
	if err != nil {
		log.Fatal(err)
	}
}
