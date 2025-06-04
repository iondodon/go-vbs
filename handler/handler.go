package handler

import (
	"net/http"
)

// Controller is a function type for HTTP handlers that return errors
type Controller func(w http.ResponseWriter, r *http.Request) error

// Handler wraps a Controller function and converts it to an http.Handler
func Handler(controller Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := controller(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
