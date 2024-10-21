package kaimono

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func writeResponse[T any](w http.ResponseWriter, code int, payload T) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, code int, err error) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ErrorResponse{Data: nil, Error: err.Error()}

	return json.NewEncoder(w).Encode(resp)
}

func logIfError(l *slog.Logger, msg string, err error) {
	if l == nil {
		return
	}

	if err == nil {
		return
	}

	l.Error(msg, "error", err)
}

type ErrorResponse struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

type Response[T any] struct {
	Data  T      `json:"data"`
	Error string `json:"error"`
}

type GetCartResponse = Response[Cart]
type GetCartByIDResponse = Response[Cart]
type CreateCartResponse = Response[Cart]
