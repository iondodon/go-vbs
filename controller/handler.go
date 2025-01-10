package controller

import (
	"net/http"
)

type Controller func(w http.ResponseWriter, r *http.Request) error

func Handler(controller Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := controller(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
