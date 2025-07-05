package httpx

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"notes-taker/internal/models"
)

func DecodeAndValidateBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, models.BaseResponse{
				Error: &models.BaseError{
					Message: "failed to decode request",
				},
			})
			return err
		}
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.BaseResponse{
			Error: &models.BaseError{
				Message: "failed to decode request",
			},
		})
		return err
	}

	validate := validator.New()
	if err := validate.Struct(dst); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.BaseResponse{
			Error: &models.BaseError{
				Message: "validation failed",
			},
		})
		return err
	}
	return nil
}
