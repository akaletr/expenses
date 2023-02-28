package main

import (
	"fmt"

	"cmd/main/main.go/internal/app"
)

func main() {
	server := app.New()
	defer func() {
		err := server.Stop()
		fmt.Println(err)
	}()

	err := server.Init()
	if err != nil {
		fmt.Println(err)
	}

	err = server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
