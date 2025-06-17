package middleware

import (
	"encoding/json"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token != "admin" {
			data := map[string]interface{}{
				"Status":  false,
				"Message": "unauthorization",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(data)
			return
		}
		next.ServeHTTP(w, r)
	})
}
