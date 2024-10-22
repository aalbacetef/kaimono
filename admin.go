package kaimono

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (svc *Service) AdminRouter(base string) *chi.Mux {
	r := chi.NewRouter()

	r.Route(base, func(r chi.Router) {
		r.Get("/{id}", svc.GetWithID)
		r.Post("/", svc.CreateWithoutSession)
		r.Put("/{id}", svc.UpdateWithID)
		r.Delete("/{id}", svc.DeleteWithID)
	})

	return r
}

// GetWithID will return the Cart if found.
//
// Errors:
//   - NotAuthorizedError if user is not authorized
//   - ErrCartNotFound if cart is not found
//
// Status Codes:
//   - 200: OK
//   - 403: Forbidden
//   - 404: Cart not found
func (svc *Service) GetWithID(w http.ResponseWriter, req *http.Request) {
	op := Operation{
		Type:     ReadOp,
		Resource: "cart",
	}

	cartID := chi.URLParam(req, "id")

	if !checkAndReportAuthorized(svc, w, req, op, cartID) {
		return
	}

	cart, err := svc.db.LookupCart(cartID)
	if errors.Is(err, ErrCartNotFound) {
		svc.json(writeError(w, http.StatusNotFound, err))
		return
	}

	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	svc.json(writeResponse(w, http.StatusOK, GetCartByIDResponse{Data: cart}))
}

// CreateWithoutSession will create an empty Cart without
// assigning it to a session.
//
// Errors:
//   - NotAuthorizedError if user is not authorized
//   - ErrCartNotFound if cart is not found
//
// Status Codes:
//   - 200: OK
//   - 403: Forbidden
//   - 404: Cart not found
func (svc *Service) CreateWithoutSession(w http.ResponseWriter, req *http.Request) {
	op := Operation{
		Type:     CreateOp,
		Resource: "cart",
	}

	if !checkAndReportAuthorized(svc, w, req, op, "") {
		return
	}

	cart, err := svc.db.CreateCart()
	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, fmt.Errorf("could not decode request: %w", err)))
		return
	}

	svc.json(writeResponse(w, http.StatusCreated, CreateCartResponse{Data: cart}))
}

// Update will update the Cart. It will override the Cart ID to ensure no accidental
// changes.
//
// Status codes:
//   - 200: Updated successfully
//   - 400: No session found for request
//   - 404: No cart found
//   - 500: unexpected error
func (svc *Service) UpdateWithID(w http.ResponseWriter, req *http.Request) {
	op := Operation{
		Type:     UpdateOp,
		Resource: "cart",
	}

	cartID := chi.URLParam(req, "id")
	if !checkAndReportAuthorized(svc, w, req, op, cartID) {
		return
	}

	payload := UpdateCartRequest{}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		svc.json(writeError(w, http.StatusBadRequest, err))
		return
	}

	foundCart, err := svc.db.LookupCart(cartID)
	if errors.Is(err, ErrCartNotFound) {
		svc.json(writeError(w, http.StatusNotFound, err))
		return
	}

	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	// NOTE: we still overwrite the payload's cart ID
	payload.Data.ID = foundCart.ID

	if err := svc.db.UpdateCart(payload.Data); err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	svc.json(writeResponse(w, http.StatusOK, UpdateCartResponse{Data: payload.Data}))
}

// Delete will delete the Cart with the supploed ID.
//
// Status codes:
//   - 204: Deleted successfully
//   - 400: No session found for request
//   - 404: No cart found
//   - 500: unexpected error
func (svc *Service) DeleteWithID(w http.ResponseWriter, req *http.Request) {
	op := Operation{
		Type:     DeleteOp,
		Resource: "cart",
	}

	cartID := chi.URLParam(req, "id")
	if !checkAndReportAuthorized(svc, w, req, op, cartID) {
		return
	}

	if err := svc.db.DeleteCart(cartID); err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func checkAndReportAuthorized(svc *Service, w http.ResponseWriter, req *http.Request, op Operation, id string) bool {
	err := svc.authorizer.AuthorizeUser(req, op, id)
	if errors.As(err, &NotAuthorizedError{}) {
		svc.json(writeError(w, http.StatusForbidden, err))
		return false
	}

	if err != nil {
		svc.json(writeError(w, http.StatusInternalServerError, err))
		return false
	}

	return true
}
