package tprpc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      *int            `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      *int   `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   Error  `json:"error"`
}

func NewResponse(id *int, result any) Response {
	return Response{
		Jsonrpc: "2.0",
		ID:      id,
		Result:  result,
	}
}

func ErrorResponse(id *int, err Error) Response {
	return Response{
		Jsonrpc: "2.0",
		ID:      id,
		Error:   err,
	}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type TrackEventParams struct {
	Type      string `json:"type"`
	Filepath  string `json:"filepath"`
	Timestamp int64  `json:"timestamp"`
}

type Server struct {
	scanner *bufio.Scanner
	ctx     context.Context
	methods map[string]func() (any, error)
}

func NewServer() Server {
	return Server{
		scanner: bufio.NewScanner(os.Stdin),
		ctx:     context.Background(),
	}
}

var (
	ParseError = Error{
		Code:    -32700,
		Message: "Parse Error",
        Data: "Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.",
	}
    InvalidRequest = Error {
        Code: -32600,
        Message: "Invalid Request",
        Data: "The JSON sent is not a valid Request object.",
    }
    InvalidParams = Error {
        Code: -32602,
        Message: "Invalid params",
        Data: "Invalid method parameter(s).",
    }
)

func MethodNotFound(method string) Error {
    return Error {
        Code: -32601,
        Message: "Method not found",
        Data: fmt.Sprintf("The method '%s' does not exist.", method),
    }
}

func MethodError(method string, err error) Error {
    return Error{
        Code: -1,
        Message: fmt.Sprintf("Method '%s' error", method),
        Data: fmt.Sprintf("Method '%s' returned the error '%v'", method, err),
    }
}


func (s Server) Write(res Response) error {
    err := json.NewEncoder(os.Stdout).Encode(res)
    return err;
}

func (s Server) ListenAndServe() {
	for s.scanner.Scan() {
		var req Request
		err := json.Unmarshal(s.scanner.Bytes(), &req)
		if err != nil {
            s.Write(ErrorResponse(nil, ParseError))
			continue
		}
		if req.Jsonrpc != "2.0" {
            s.Write(ErrorResponse(nil, InvalidRequest))
			continue
		}
		if req.ID == nil {
            s.Write(ErrorResponse(nil, InvalidRequest))
			continue
		}
		method, ok := s.methods[req.Method]
		if !ok {
            s.Write(ErrorResponse(req.ID, MethodNotFound(req.Method)))
            continue
		}
        res, err := method()
        if err != nil {
            s.Write(ErrorResponse(req.ID, MethodError(req.Method, err)))
            continue
        }
        s.Write(NewResponse(req.ID, res))
	}
}
