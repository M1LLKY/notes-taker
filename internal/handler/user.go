package handler

import (
	"github.com/go-chi/render"
	"net/http"
	"notes-taker/internal/errcodes"
	"notes-taker/internal/httpx"
	"notes-taker/internal/service"
)

func (h *Handler) PostSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var signIn SignInRequest
	if err := httpx.DecodeAndValidateBody(w, r, &signIn); err != nil {
		return
	}
	input := service.SignIn{
		Username: signIn.Username,
		Password: signIn.Password,
	}
	response, err := h.deps.UserService.SignIn(ctx, input)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var signIn SignInRequest
	if err := httpx.DecodeAndValidateBody(w, r, &signIn); err != nil {
		return
	}
	input := service.SignIn{
		Username: signIn.Username,
		Password: signIn.Password,
	}
	response, err := h.deps.UserService.Login(ctx, input)
	if err != nil {
		errcodes.SendErrorJSON(w, r, http.StatusInternalServerError, err)
		return
	}
	render.JSON(w, r, response)
}
