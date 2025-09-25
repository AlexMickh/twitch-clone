package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type HttpError struct {
	error
	msg    string
	status int
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorWrapper(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			httpErr, ok := err.(HttpError)
			if ok {
				render.Status(r, httpErr.status)
				render.JSON(w, r, ErrorResponse{Error: httpErr.msg})
			} else {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, ErrorResponse{Error: err.Error()})
			}
		}
	}
}

func Error(msg string, status int) HttpError {
	return HttpError{
		msg:    msg,
		status: status,
	}
}
