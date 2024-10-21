package kaimono

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (svc *Service) Router(base string) *chi.Mux {
	r := chi.NewRouter()

	r.Route(base, func(r chi.Router) {
		r.Get("/", svc.Get)
		r.Post("/", svc.Create)
		r.Put("/", svc.Update)
		r.Delete("/", svc.Delete)
	})

	return r
}

/*************************************
*
*         Standard routes
*
**************************************/

// Get will return the Cart associated to the current user's session.
//
// Status codes:
//   - 200: OK
//   - 400: No session found for request
//   - 404: No cart found for session
//   - 500: unexpected error
func (svc *Service) Get(w http.ResponseWriter, req *http.Request) {
	usrCtx, err := svc.usrCtxFetcher.GetUserContext(req)
	if errors.Is(err, ErrSessionNotFound) {
		svc.json(
			writeError(w, http.StatusBadRequest, ErrSessionNotFound),
		)

		return
	}

	if err != nil {
		svc.json(
			writeError(w, http.StatusInternalServerError, err),
		)
		return
	}

	cart, err := svc.db.LookupCartForSession(usrCtx.SessionToken)
	if errors.Is(err, ErrCartNotFound) {
		svc.json(writeError(w, http.StatusNotFound, ErrCartNotFound))
		return
	}

	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	svc.json(writeResponse(w, GetCartResponse{Data: cart}))
}

// Create will create a new Cart for the current session.
//
// Status codes:
//   - 201: Created successfully
//   - 400: No session found for request
//   - 409: Cart already exists
//   - 500: unexpected error
func (svc *Service) Create(w http.ResponseWriter, req *http.Request) {
	usrCtx, err := svc.usrCtxFetcher.GetUserContext(req)
	if errors.Is(err, ErrSessionNotFound) {
		svc.json(
			writeError(w, http.StatusBadRequest, ErrSessionNotFound),
		)

		return
	}

	if err != nil {
		svc.json(
			writeError(w, http.StatusInternalServerError, err),
		)
		return
	}

	cart, err := svc.db.CreateCart(usrCtx.SessionToken)
	if errors.Is(err, ErrAlreadyExists) {
		svc.json(writeError(w, http.StatusConflict, err))
		return
	}

	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	// @TODO: add Location header
	svc.json(writeResponse(w, CreateCartResponse{Data: cart}))
}

// Update will create the Cart for the current session. It will reject
// the Cart if the ID suplied does not match.
//
// Status codes:
//   - 201: Created successfully
//   - 400: No session found for request
//   - 403: Cart ID is not the ID matching this session's Cart
//   - 404: No cart found for this session
//   - 500: unexpected error
func (svc *Service) Update(w http.ResponseWriter, req *http.Request) {
	usrCtx, err := svc.usrCtxFetcher.GetUserContext(req)
	if errors.Is(err, ErrSessionNotFound) {
		svc.json(
			writeError(w, http.StatusBadRequest, ErrSessionNotFound),
		)

		return
	}

	if err != nil {
		svc.json(
			writeError(w, http.StatusInternalServerError, err),
		)
		return
	}

	foundCart, err := svc.db.LookupCartForSession(usrCtx.SessionToken)
	if errors.Is(err, ErrCartNotFound) {
		svc.json(writeError(w, http.StatusNotFound, ErrCartNotFound))
		return
	}

	payload := UpdateCartRequest{}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		svc.json(writeError(w, http.StatusBadRequest, fmt.Errorf("could not decode request: %w", err)))
		return
	}

	if payload.Data.ID != foundCart.ID {
		svc.json(writeError(w, http.StatusForbidden, ErrInvalidID))
	}

	if err := svc.db.UpdateCart(payload.Data); err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, fmt.Errorf("update failed: %w", err)))
		return
	}

}

func (svc *Service) Delete(w http.ResponseWriter, req *http.Request) {
	usrCtx, err := svc.usrCtxFetcher.GetUserContext(req)
	if errors.Is(err, ErrSessionNotFound) {
		svc.json(
			writeError(w, http.StatusBadRequest, ErrSessionNotFound),
		)

		return
	}

	if err != nil {
		svc.json(
			writeError(w, http.StatusInternalServerError, err),
		)
		return
	}

	foundCart, err := svc.db.LookupCartForSession(usrCtx.SessionToken)
	if errors.Is(err, ErrCartNotFound) {
		svc.json(writeError(w, http.StatusNotFound, ErrCartNotFound))
		return
	}

	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	if err := svc.db.DeleteCart(foundCart.ID); err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
