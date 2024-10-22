package kaimono

import (
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
	data     map[string]Cart
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
		data:     make(map[string]Cart),
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

func (mock *mockBackend) CreateCart(sessionToken string) (Cart, error) {
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

	cart, found := mock.data[sessionToken]
	if !found {
		return cart, ErrCartNotFound
	}

	return cart, nil
}

func (mock *mockBackend) DeleteCart(cartID string) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	for key, cart := range mock.data {
		if cart.ID == cartID {
			delete(mock.data, key)
			return nil
		}
	}

	return ErrCartNotFound
}

func (mock *mockBackend) UpdateCart(cart Cart) error {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	for key, _cart := range mock.data {
		if _cart.ID == cart.ID {
			mock.data[key] = cart
			return nil
		}
	}

	return ErrCartNotFound
}

func (mock *mockBackend) LookupCart(cartID string) (Cart, error) {
	mock.mu.Lock()
	defer mock.mu.Unlock()

	for _, cart := range mock.data {
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

	cart, found := mock.data[sessionToken]
	if !found {
		return cart, ErrCartNotFound
	}

	return cart, nil
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
