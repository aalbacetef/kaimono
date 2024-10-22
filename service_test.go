package kaimono

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

const testCookieName = "test-cookie"

// mockBackend implements the main service interfaces
type mockBackend struct {
	users    []UserContext
	sessions []string
	data     map[string]int
	carts    []Cart
	mu       sync.Mutex
}

func mkEmptyTestCart() Cart {
	return Cart{
		ID:        uuid.New().String(),
		Items:     []CartItem{},
		Discounts: []Discount{},
	}
}

func newMockBackend() *mockBackend {
	sessions := []string{
		"logged-in-session",
		"logged-in-admin-session",
		"anonymous-session",
	}
	users := []UserContext{
		{UserID: "test-user", SessionToken: "logged-in-session"},
		{UserID: "test-admin-user", SessionToken: "logged-in-admin-session"},
	}

	mock := &mockBackend{
		carts:    []Cart{},
		data:     make(map[string]int),
		sessions: sessions,
		users:    users,
	}

	return mock
}

func (mock *mockBackend) GetUserContext(req *http.Request) (UserContext, error) {
	cookie, err := req.Cookie(testCookieName)
	if err != nil {
		return UserContext{}, fmt.Errorf("could not get cookie: %w", err)
	}

	usrCtx := UserContext{
		SessionToken: cookie.Value,
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	for _, u := range mock.users {
		if u.SessionToken == usrCtx.SessionToken {
			return u, nil
		}
	}

	return usrCtx, nil
}

func (mock *mockBackend) CreateCartForSession(sessionToken string) (Cart, error) {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	index := -1
	for k, s := range mock.sessions {
		if s == sessionToken {
			index = k
			break
		}
	}

	if index == -1 {
		return Cart{}, ErrSessionNotFound
	}

	cartIndex, found := mock.data[sessionToken]
	if found {
		return mock.carts[cartIndex], ErrAlreadyExists
	}

	cart := mkEmptyTestCart()
	mock.carts = append(mock.carts, cart)
	cartIndex = len(mock.carts) - 1

	mock.data[sessionToken] = cartIndex

	return cart, nil
}

func (mock *mockBackend) CreateCart() (Cart, error) {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	cart := mkEmptyTestCart()
	mock.carts = append(mock.carts, cart)

	return cart, nil
}

func (mock *mockBackend) DeleteCart(cartID string) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	index := -1
	for k, cart := range mock.carts {
		if cart.ID == cartID {
			index = k
			break
		}
	}

	if index == -1 {
		return ErrCartNotFound
	}

	// delete entry for session, if found
	for key, indx := range mock.data {
		cart := mock.carts[indx]
		if cart.ID == cartID {
			delete(mock.data, key)
			break
		}
	}

	// remove from carts
	mock.carts = removeAt[Cart](mock.carts, index)

	return nil
}

func (mock *mockBackend) UpdateCart(cart Cart) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	for k, c := range mock.carts {
		if c.ID == cart.ID {
			mock.carts[k] = cart
			return nil
		}
	}

	return ErrCartNotFound
}

func (mock *mockBackend) LookupCart(cartID string) (Cart, error) {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	for _, cart := range mock.carts {
		if cart.ID == cartID {
			return cart, nil
		}
	}

	return Cart{}, ErrCartNotFound
}

func (mock *mockBackend) LookupCartForSession(sessionToken string) (Cart, error) {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	index := -1
	for k, s := range mock.sessions {
		if s == sessionToken {
			index = k
			break
		}
	}

	if index == -1 {
		return Cart{}, ErrSessionNotFound
	}

	cartIndex, found := mock.data[sessionToken]
	if !found {
		return Cart{}, ErrCartNotFound
	}

	return mock.carts[cartIndex], nil
}

func (mock *mockBackend) AssignCartToSession(string, string) error {
	return errors.New("not implemented")
}

func (mock *mockBackend) AuthorizeUser(req *http.Request, op Operation, id string) error {
	usrCtx, err := mock.GetUserContext(req)
	if err != nil {
		return err
	}

	if usrCtx.UserID == "test-admin-user" {
		return nil
	}

	return NotAuthorizedError{Operation: op, ID: id}
}

func removeAt[T any](arr []T, index int) []T {
	return append(arr[:index-1], arr[index:]...)
}
