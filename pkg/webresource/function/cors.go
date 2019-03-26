package function

import (
	"net/http"
	"strconv"
	"strings"
)

func handlePreFlight(w http.ResponseWriter, r *http.Request) (bool, http.ResponseWriter) {
	origin := r.Header.Get("Origin")

	if origin == "" {
		origin = "*"
	}

	headers := w.Header()
	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Add("Vary", "Origin")

	if r.Method == "OPTIONS" {
		allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		allowedHeaders := []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}
		exposedHeaders := []string{"Link"}

		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
		headers.Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
		headers.Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		headers.Set("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ", "))
		headers.Set("Access-Control-Max-Age", strconv.Itoa(3600))

		return true, w
	}

	return false, w
}
