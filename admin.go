package kaimono

import (
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

/*************************************
*
*         Admin routes
*
**************************************/

func (svc *Service) GetWithID(w http.ResponseWriter, req *http.Request) {}

func (svc *Service) CreateWithoutSession(w http.ResponseWriter, req *http.Request) {}

func (svc *Service) UpdateWithID(w http.ResponseWriter, req *http.Request) {}

func (svc *Service) DeleteWithID(w http.ResponseWriter, req *http.Request) {}
