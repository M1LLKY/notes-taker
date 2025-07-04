package common

import (
	"bytes"
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter
	Status int
	Size   int64
	Body   *bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Body:           &bytes.Buffer{},
	}
}

func (r *ResponseRecorder) WriteHeader(code int) {
	r.Status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	if r.Body != nil {
		r.Body.Write(b)
	}
	size, err := r.ResponseWriter.Write(b)
	r.Size += int64(size)
	return size, err
}
