package app

import "cmd/main/main.go/internal/config"

type app struct {
}

func NewApp(cfg config.Config) App {
	return &app{}
}

func (app *app) Init() error {
	return nil
}

func (app *app) Start() error {
	return nil
}

func (app *app) Stop() error {
	return nil
}
