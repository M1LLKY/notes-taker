package errcodes

import (
	"github.com/go-chi/render"
	"net/http"
	"notes-taker/internal/models"
)

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpCode int, err error) {
	errorResponse := models.BaseResponse{
		Error: &models.BaseError{Message: err.Error()},
	}
	render.Status(r, httpCode)
	render.JSON(w, r, errorResponse)
}
