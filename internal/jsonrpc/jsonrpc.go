package jsonrpc

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	ID     string          `json:"id"`
	Error  string          `json:"error"`
	Result json.RawMessage `json:"result"`
}

type Options struct {
	Conn   *gorm.DB
	Params json.RawMessage
}

type Method func(options Options) (json.RawMessage, error)
