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

	r.Route("/", func(r chi.Router) {
		r.Post("/sign-in", h.PostSignIn)
		r.Post("/login", h.PostLogin)
	})

	r.With(h.jwtAuthMiddleware).Group(func(r chi.Router) {
		r.Post("/users/{id}/notes", h.PostCreateNote)
		r.Get("/users/{id}/notes", h.GetNote)
		r.Get("/users/{id}/notes/{note_id}", h.GetNoteByID)
		r.Put("/users/{id}/notes/{note_id}", h.PutNoteByID) // Изменение заметки
		r.Delete("/users/{id}/notes/{note_id}", h.DeleteNoteByID)
	})

	return r
}
