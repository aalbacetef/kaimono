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

	// DeleteCartForSession will delete the Cart for the given session.
	//
	// If no matching session is found, it will return ErrSessionNotFound.
	// If no Cart for that session exists, it will return ErrCartNotFound.
	DeleteCartForSession(sessionToken string) error

	// DeleteCart will delete the Cart matching the ID. It doesn't check
	// for permissions and shouldn't be used except for admin purposes.
	DeleteCart(cartID string) error

	UpdateCart(cart Cart) error
	LoadCart(cartID string)

	LookupCartForSession(sessionToken string) (Cart, error)
}

func DeleteCartForSession(db DB, sessionToken string) error {
	cart, err := db.LookupCartForSession(sessionToken)
	if err != nil {
		return err
	}

	return db.DeleteCart(cart.ID)
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
