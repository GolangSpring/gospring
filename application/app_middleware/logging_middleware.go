package appMiddleware

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (wrapper *ResponseWriterWrapper) WriteHeader(statusCode int) {
	wrapper.StatusCode = statusCode
	wrapper.ResponseWriter.WriteHeader(statusCode)
}

type (
	SensitiveContext struct {
		Password   string `json:"password"`
		Pwd        string `json:"pwd"`
		Token      string `json:"token"`
		SSN        string `json:"ssn"`         // Social Security Number
		CreditCard string `json:"credit_card"` // Credit card number
		Email      string `json:"email"`       // User email
		Phone      string `json:"phone"`       // User phone number
		Address    string `json:"address"`     // Home address
		DOB        string `json:"dob"`         // Date of birth
		APIKey     string `json:"api_key"`     // API key
	}
)

func (ctx *SensitiveContext) IsSensitive() bool {
	return len(ctx.Password) > 0 ||
		len(ctx.Pwd) > 0 ||
		len(ctx.Token) > 0 ||
		len(ctx.SSN) > 0 ||
		len(ctx.CreditCard) > 0 ||
		len(ctx.Email) > 0 ||
		len(ctx.Phone) > 0 ||
		len(ctx.Address) > 0 ||
		len(ctx.DOB) > 0 ||
		len(ctx.APIKey) > 0
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)
		// Log the request details

		statusCode := wrapper.StatusCode
		logger := log.Info() // Default log level

		switch {
		case statusCode >= 500:
			logger = log.Error() // Server errors
		case statusCode >= 400:
			logger = log.Warn() // Client errors
		}

		logger.
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Int("status", statusCode).
			Str("user_agent", r.UserAgent()).
			Dur("duration", time.Since(start)).
			Msg("Request processed")
	})

}
