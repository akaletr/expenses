package app

import "cmd/main/main.go/internal/jsonrpc"

type App interface {
	Init() error
	Start() error
	Stop() error

	Register(name string, method jsonrpc.Method)
}
