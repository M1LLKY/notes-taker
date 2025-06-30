package create_user

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	resp "notes-taker/internal/lib/api/response"
	"notes-taker/internal/storage"
)

type Request struct {
	Username string `json:"username,omitempty"`
}

type Response struct {
	resp.Response
	Username string `json:"username,omitempty"`
}

type UserCreator interface {
	CreateUser(username string) error
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.create-user.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "empty request",
			})

			return
		}
		if err != nil {
			log.Error("failed to decode request body")

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "failed to decode request",
			})

			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		username := req.Username
		if username == "" {
			log.Error("empty username")

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "empty username",
			})

			return
		}

		err = userCreator.CreateUser(req.Username)
		if errors.Is(err, storage.UserExists) {
			log.Info("user already exists")

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "user already exists",
			})

			return
		}
		if err != nil {
			slog.Error("failed to add user:", err)

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "failed to add user",
			})

			return
		}

		log.Info("user added")

		responseOK(w, r, username)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, username string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Username: username,
	})
}
