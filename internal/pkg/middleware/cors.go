package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		allowedOrigins := map[string]bool{
			"http://localhost":      true,
			"http://localhost:8080": true,
			"http://localhost:4200": true,
			"http://127.0.0.1":      true,
			"http://127.0.0.1:8080": true,
			"http://127.0.0.1:4200": true,
		}

		header := w.Header()

		if allowedOrigins[origin] {
			header.Set("Access-Control-Allow-Origin", origin)
			header.Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			header.Set("Access-Control-Allow-Methods", "POST, GET, DELETE, HEAD, PATCH, PUT")
			header.Set("Access-Control-Allow-Headers", "authorization,content-type,content-length")
			header.Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return

		}
		next.ServeHTTP(w, r)
	})
}
