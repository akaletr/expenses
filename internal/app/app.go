package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"cmd/main/main.go/internal/auth"
	"cmd/main/main.go/internal/config"
	"cmd/main/main.go/internal/entity"
	"cmd/main/main.go/internal/storage"

	"github.com/go-chi/chi/v5"
)

type app struct {
	auth    auth.Auth
	storage storage.Storage
	cfg     config.Config
	server  http.Server
}

func NewApp(cfg config.Config) (App, error) {
	db, err := storage.New()
	if err != nil {
		return nil, err
	}

	return &app{
		storage: db,
		cfg:     cfg,
		auth:    auth.New(cfg.SecretKey),
	}, nil
}

func (app *app) Init() error {
	err := app.storage.Start(app.cfg.Database)
	if err != nil {
		return err
	}

	err = app.storage.Provide(
		&entity.Wallet{},
		&entity.User{},
		&entity.Category{},
		&entity.Event{},
	)
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

	router := chi.NewRouter()
	router.Use(app.auth.CookieHandler)

	router.Get("/category/{id}", app.getCategory)
	router.Get("/categories", app.getCategories)
	router.Put("/category", app.putCategory)

	router.Get("/event/{id}", app.getEvent)
	router.Get("/events", app.getEvents)
	router.Put("/event", app.putEvent)

	app.server.Handler = router
	return nil
}

func (app *app) Start() error {
	return app.server.ListenAndServe()
}

func (app *app) Stop() error {
	app.storage.Stop()
	return nil
}

func (app *app) getCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	c := entity.Category{}

	err = c.Get(app.storage.GetDB(), uint(idInt))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (app *app) getCategories(w http.ResponseWriter, r *http.Request) {
	cc, err := r.Cookie("user")
	if err != nil {
		cc = &http.Cookie{}
	}

	id, err := app.auth.GetID(cc)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("________", id)
	time.Sleep(time.Second / 2)
	var c []entity.Category
	app.storage.GetDB().Where("user_id = ?", id).Find(&c)

	data, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (app *app) putCategory(w http.ResponseWriter, r *http.Request) {
	cc, err := r.Cookie("user")
	if err != nil {
		cc = &http.Cookie{}
	}

	id, err := app.auth.GetID(cc)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	c := entity.Category{}

	i, _ := strconv.Atoi(id)

	c.UserID = uint(i)

	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = c.Put(app.storage.GetDB())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("ok"))
}

func (app *app) getEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	e := entity.Event{}

	err = e.Get(app.storage.GetDB(), uint(idInt))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	data, err := json.Marshal(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (app *app) putEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	e := entity.Event{}

	err = json.Unmarshal(body, &e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println(e)
	err = e.Put(app.storage.GetDB())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("ok"))
}

func (app *app) getEvents(w http.ResponseWriter, r *http.Request) {
	var c []entity.Event
	app.storage.GetDB().Order("updated_at DESC").Find(&c)

	data, err := json.Marshal(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
