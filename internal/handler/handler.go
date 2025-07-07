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
		r.Post("/users/notes", h.PostCreateNote)
		r.Get("/users/notes", h.GetAllNotes)
		r.Get("/users/notes/{note_id}", h.GetNoteByID)
		r.Put("/users/notes/{note_id}", h.PutUpdateNoteByID) // Изменение заметки
		r.Delete("/users/notes/{note_id}", h.DelDeleteNoteByID)
	})

	return r
}
