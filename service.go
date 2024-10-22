package kaimono

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidID       = errors.New("invalid ID")
)

type Service struct {
	authorizer    Authorizer
	db            DB
	usrCtxFetcher UserContextFetcher
	logger        *slog.Logger
}

func NewService(db DB, usrCtxFetcher UserContextFetcher, authorizer Authorizer, logger *slog.Logger) (*Service, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	return &Service{
		authorizer:    authorizer,
		db:            db,
		usrCtxFetcher: usrCtxFetcher,
		logger:        logger,
	}, nil
}

func (svc *Service) json(err error) {
	logIfError(svc.logger, "write response", err)
}

type Authorizer interface {

	// AuthorizeUser will determine if the user (retrieved from the request)
	// can perform the given operation on the specified resource.
	AuthorizeUser(req *http.Request, op Operation, id string) error
}

type NotAuthorizedError struct {
	Operation Operation `json:"operation"`
	ID        string    `json:"id"`
}

func (e NotAuthorizedError) Error() string {
	return fmt.Sprintf("could not authorize user for op=(%+v) with resource ID:(%s)", e.Operation, e.ID)
}

type Operation struct {
	Resource string        `json:"resource"`
	Type     OperationType `json:"type"`
}

type OperationType string

const (
	CreateOp OperationType = "create"
	ReadOp   OperationType = "read"
	UpdateOp OperationType = "update"
	DeleteOp OperationType = "delete"
)

type DB interface {
	// CreateCartForSession will instantiate a brand new empty Cart for the session.
	//
	// If no matching session is found it will return ErrSessionNotFound.
	// If a Cart already exists for that session, it will return ErrAlreadyExists.
	CreateCartForSession(sessionToken string) (Cart, error)

	// CreateCart will create a Cart without assigning it to a session.
	CreateCart() (Cart, error)

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
	LookupCart(cartID string) (Cart, error)

	// LookupCart will find the Cart for this session.
	//
	// If no matching session is found, it will return ErrSessionNotFound.
	// If no cart could be found, it will return ErrCartNotFound.
	LookupCartForSession(sessionToken string) (Cart, error)

	// AssignCartToSession will assign the cart specified by ID to the given
	// session.
	//
	// If no matching session is found, it will return ErrSessionNotFound.
	// If no cart could be found, it will return ErrCartNotFound.
	AssignCartToSession(cartID, sessionToken string) error
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
