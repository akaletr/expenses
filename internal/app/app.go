package app

import (
	"cmd/main/main.go/internal/entity/category"
	"cmd/main/main.go/internal/entity/event"
	"cmd/main/main.go/internal/entity/user"
	"cmd/main/main.go/internal/entity/wallet"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"cmd/main/main.go/internal/auth"
	"cmd/main/main.go/internal/config"
	"cmd/main/main.go/internal/jsonrpc"
	"cmd/main/main.go/internal/storage"

	"github.com/go-chi/chi/v5"
)

type app struct {
	auth     auth.Auth
	storage  storage.Storage
	cfg      config.Config
	server   http.Server
	handlers map[string]jsonrpc.Method
	mu       sync.RWMutex
}

func NewApp(cfg config.Config) (App, error) {
	db, err := storage.New()
	if err != nil {
		return nil, err
	}

	return &app{
		storage:  db,
		cfg:      cfg,
		auth:     auth.New(cfg.SecretKey),
		handlers: map[string]jsonrpc.Method{},
	}, nil
}

func (app *app) Register(name string, method jsonrpc.Method) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.handlers[strings.ToLower(name)] = method
}

func (app *app) Init() error {
	err := app.storage.Start(app.cfg.Database)
	if err != nil {
		return err
	}

	err = app.storage.Register(
		&wallet.Wallet{},
		&user.User{},
		&category.Category{},
		&event.Event{},
	)
	if err != nil {
		return err
	}

	app.Register("getWallet", wallet.GetWallet)
	app.Register("getUser", wallet.GetWallet)

	app.Register("category.get", category.Get)
	app.Register("category.getMany", category.GetMany)
	app.Register("category.create", category.Create)

	app.Register("getEvent", wallet.GetWallet)

	app.server = http.Server{
		Addr:              fmt.Sprintf(":%s", app.cfg.ServerPort),
		Handler:           nil,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
	}

	router := chi.NewRouter()
	router.Use(app.auth.CookieHandler)
	router.Post("/", app.handleRequest)

	router.Get("/category/{id}", app.getCategory)
	router.Get("/categories", app.getCategories)
	router.Put("/category", app.putCategory)

	router.Get("/event/{id}", app.getEvent)
	router.Get("/events", app.getEvents)
	router.Put("/event", app.putEvent)

	router.Post("/login", app.login)

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
	c := category.Category{}

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
	var c []category.Category
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

	c := category.Category{}

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
	e := event.Event{}

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

	e := event.Event{}

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
	var c []event.Event
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

func (app *app) login(w http.ResponseWriter, r *http.Request) {
	var user user.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = user.Put(app.storage.GetDB())
	http.Redirect(w, r, "http://localhost:8080/events", http.StatusSeeOther)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}

func (app *app) getMethod(name string) (jsonrpc.Method, error) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	method, ok := app.handlers[name]
	if ok {
		return method, nil
	}

	return nil, errors.New(fmt.Sprintf("'%s' not found", name))
}

func (app *app) handleRequest(w http.ResponseWriter, r *http.Request) {
	response := jsonrpc.Response{}
	defer func() {
		data, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Error = err.Error()
		return
	}

	var request jsonrpc.Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Error = err.Error()
		return
	}

	method, err := app.getMethod(request.Method)
	if err != nil {
		response.Error = err.Error()
		return
	}

	options := jsonrpc.Options{
		Conn:   app.storage.GetDB(),
		Params: request.Params,
	}

	result, err := method(options)
	if err != nil {
		response.Error = err.Error()
		return
	}

	response.ID = request.ID
	response.Result = result
}
