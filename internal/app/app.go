package app

import (
	"cmd/main/main.go/internal/gziper"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func (app *app) register(name string, method jsonrpc.Method) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.handlers[strings.ToLower(name)] = method
}

func (app *app) init() error {
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

	app.register("category.get", category.Get)
	app.register("category.getMany", category.GetMany)
	app.register("category.create", category.Create)
	app.register("category.delete", category.Delete)

	app.register("event.get", event.Get)
	app.register("event.getMany", event.GetMany)
	app.register("event.create", event.Create)
	app.register("event.delete", event.Delete)

	app.register("wallet.get", wallet.Get)
	app.register("wallet.getMany", wallet.GetMany)
	app.register("wallet.create", wallet.Create)
	app.register("wallet.delete", wallet.Delete)

	app.register("user.get", user.Get)
	app.register("user.getMany", user.GetMany)
	app.register("user.create", user.Create)
	app.register("user.delete", user.Delete)

	app.server = http.Server{
		Addr:              fmt.Sprintf(":%s", app.cfg.ServerPort),
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
	}

	router := chi.NewRouter()
	router.Use(gziper.GzipCompression)

	router.Post("/v1", app.handleRequest)
	router.Post("/v1/auth", app.authn)

	app.server.Handler = router
	return nil
}

func (app *app) Start() error {
	err := app.init()
	if err != nil {
		return err
	}

	fmt.Println("Server started at port:", app.cfg.ServerPort)
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

	fmt.Println("sdklfjalsdkjf")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", " GET, PUT, POST, DELETE, OPTIONS")

	defer func() {
		data, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}()

	idCookie, err := r.Cookie("user")
	if err != nil {
		response.Error = err.Error()
		return
	}
	id, err := strconv.Atoi(idCookie.Value)
	if err != nil {
		response.Error = err.Error()
		return
	}

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
		UserId: uint(id),
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

func (app *app) authn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", " GET, PUT, POST, DELETE, OPTIONS")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	conn := app.storage.GetDB()
	var candidate user.User
	err = json.Unmarshal(body, &candidate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx := conn.Create(&candidate)
	if tx.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "user",
		Value: "8",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
}
