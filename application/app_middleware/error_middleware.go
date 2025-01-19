package appMiddleware

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime/debug"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				stacks := string(debug.Stack())
				fmt.Println(stacks)

				// Log the stack trace
				fmt.Println("Error:", rec)

				// Respond with a custom error response
				w.WriteHeader(http.StatusBadRequest)
				jsonResponse, _ := json.Marshal(map[string]string{
					"message": fmt.Sprintf("%v", rec),
				})
				_, err := w.Write(jsonResponse)
				if err != nil {
					log.Error().Msgf("Failed to write response: %v", err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
