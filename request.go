package kaimono

type Request[T any] struct {
	Data T `json:"data"`
}

type UpdateCartRequest = Request[Cart]
