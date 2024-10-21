package kaimono

import (
	"errors"
	"log/slog"
	"net/http"
)

type Service struct {
	db            DB
	usrCtxFetcher UserContextFetcher
	l             *slog.Logger
}

func (svc *Service) json(err error) {
	logIfError(svc.l, "write response", err)
}

type DB interface {
	// CreateCart will instantiate a brand new empty Cart for the session.
	//
	// If no matching session is found it will return ErrSessionNotFound.
	// If a Cart already exists for that session, it will return ErrAlreadyExists.
	CreateCart(sessionToken string) (Cart, error)

	// DeleteCart will delete the Cart matching the ID. It doesn't check
	// for permissions and should only be called after user has been authorized.
	//
	// If no Cart could be found, it will return ErrCartNotFound.
	DeleteCart(cartID string) error

	// UpdateCart will update the cart matching the cart.ID field. It doesn't check
	// for permissions and should only be called after user has been authorized.
	//
	// If no Cart could be found, it will return ErrCartNotFound.
	UpdateCart(cart Cart) error

	// LookupCart will find the Cart matching the ID. It doesn't check
	// for permissions and should only be called after user has been authorized.
	//
	// If no cart could be found, it will return ErrCartNotFound.
	LookupCart(cartID string)

	// LookupCart will find the Cart for this session.
	//
	// If no matching session is found, it will return ErrSessionNotFound.
	// If no cart could be found, it will return ErrCartNotFound.
	LookupCartForSession(sessionToken string) (Cart, error)
}

// UserContextFetcher encapsulates functionality
// for fetching session tokens and user IDs from a
// request.
// Returns ErrSessionNotFound if no session could be
// found.
type UserContextFetcher interface {
	GetUserContext(req *http.Request) (UserContext, error)
}

type UserContext struct {
	UserID       string
	SessionToken string
}

func (u UserContext) IsLoggedIn() bool {
	return u.UserID != ""
}

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidID       = errors.New("invalid ID")
)
