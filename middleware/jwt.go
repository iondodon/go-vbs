package middleware

import (
	"fmt"
	"net/http"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("checking jwt")

		next.ServeHTTP(w, r)

		fmt.Println("after checking jwt")
	})
}
