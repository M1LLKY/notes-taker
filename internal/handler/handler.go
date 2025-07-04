package handler

import (
	"github.com/go-chi/chi/v5"
	"notes-taker/internal/dependencies"
)

type Handler struct {
	deps *dependencies.Dependencies
}

func New(deps *dependencies.Dependencies) *Handler {
	return &Handler{
		deps: deps,
	}
}

func (h *Handler) GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(h.corsMiddleware())
	r.Use(h.commonWare)
	r.Use(h.handlerLogger)

	return r
}
