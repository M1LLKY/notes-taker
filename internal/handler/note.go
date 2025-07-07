package handler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"notes-taker/internal/auth"
	"notes-taker/internal/errcodes"
	"notes-taker/internal/httpx"
	"notes-taker/internal/service"
	"strconv"
)

func (h *Handler) PostCreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
		return
	}
	var request CreateNoteRequest
	if err := httpx.DecodeAndValidateBody(w, r, &request); err != nil {
		return
	}
	input := &service.CreateNote{
		Title:   request.Title,
		Content: request.Content,
		UserID:  userID,
	}
	response, err := h.deps.NoteService.CreateNote(ctx, input)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
		return
	}
	response, err := h.deps.NoteService.GetAllNotes(ctx, userID)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}

func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
		return
	}
	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("неверный формат идентификатора задачи"))
		return
	}
	response, err := h.deps.NoteService.GetNoteByID(ctx, userID, noteID)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}

func (h *Handler) PutUpdateNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
		return
	}
	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("неверный формат идентификатора задачи"))
		return
	}
	var request UpdateNoteRequest
	if err := httpx.DecodeAndValidateBody(w, r, &request); err != nil {
		return
	}
	input := &service.UpdateNote{
		ID:      noteID,
		Title:   request.Title,
		Content: request.Content,
		UserID:  userID,
	}
	response, err := h.deps.NoteService.UpdateNoteByID(ctx, input)
	if err != nil {
		if errors.Is(err, service.ErrNoteNotFound) {
			errcodes.SendErrorJSON(w, r, http.StatusNotFound, err)
			return
		}
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}

func (h *Handler) DelDeleteNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusUnauthorized, errors.New("пользователь не авторизован"))
		return
	}
	noteIDStr := chi.URLParam(r, "note_id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("неверный формат идентификатора задачи"))
		return
	}
	response, err := h.deps.NoteService.DeleteNoteByID(ctx, userID, noteID)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}
