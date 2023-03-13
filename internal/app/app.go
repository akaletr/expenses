package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"cmd/main/main.go/internal/auth"
	"cmd/main/main.go/internal/config"
	"cmd/main/main.go/internal/entity/category"
	"cmd/main/main.go/internal/entity/event"
	"cmd/main/main.go/internal/entity/user"
	"cmd/main/main.go/internal/entity/wallet"
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

func NewApp(cfg config.Config) App {
	return &app{
		cfg:      cfg,
		auth:     auth.New(cfg.SecretKey),
		handlers: map[string]jsonrpc.Method{},
	}
}

func (app *app) Register(name string, method jsonrpc.Method) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.handlers[strings.ToLower(name)] = method
}

func (app *app) Init() error {
	db, err := storage.New()
	if err != nil {
		return err
	}
	app.storage = db

	err = app.storage.Start(app.cfg.Database)
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
	app.Register("category.delete", category.Delete)

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
