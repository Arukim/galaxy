package core

import (
	"encoding/json"
)

// CommandJSON is model for network commands
type CommandJSON struct {
	Cmd  string           `json:"cmd"`
	Data *json.RawMessage `json:"data"`
}

// Command is internal command structure, command usually are sent from
// server to client
type Command struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

// ResponseJSON is sent over the web socket to client
type ResponseJSON struct {
	Cmd    string  `json:"cmd"`
	Result *Result `json:"result"`
}

// Result is internal response from command handler
type Result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// NewSuccessResult creates success result
func NewSuccessResult(data interface{}) *Result {
	r := Result{
		Status: 200,
		Data:   data,
	}
	return &r
}

// NewErrorResult creates result with error
func NewErrorResult(data interface{}) *Result {
	r := Result{
		Status: 400,
		Data:   data,
	}
	return &r
}

// CommandHandler handles one command
type CommandHandler struct {
	Name   string
	Handle func(Data *json.RawMessage, p *Player) *Result
}
