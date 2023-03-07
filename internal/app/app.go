package app

import (
	"fmt"
	"net/http"
	"time"

	"cmd/main/main.go/internal/config"
	"cmd/main/main.go/internal/entity/category"
	"cmd/main/main.go/internal/entity/event"
	"cmd/main/main.go/internal/storage"
)

type app struct {
	storage storage.Storage
	cfg     config.Config
	server  http.Server
}

func NewApp(cfg config.Config) (App, error) {
	db, err := storage.New()
	if err != nil {
		return nil, nil
	}

	return &app{
		storage: db,
		cfg:     cfg,
	}, nil
}

func (app *app) Init() error {
	err := app.storage.Start(app.cfg.Database)
	if err != nil {
		return err
	}

	err = app.storage.Provide(&category.Category{}, &event.Event{})
	if err != nil {
		return err
	}

	app.server = http.Server{
		Addr:              fmt.Sprintf(":%s", app.cfg.ServerPort),
		Handler:           nil,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
	}

	return nil
}

func (app *app) Start() error {
	return app.server.ListenAndServe()
}

func (app *app) Stop() error {
	app.storage.Stop()
	return nil
}
