package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/bangueco/auction-api/internal/handlers/helper"
	"github.com/bangueco/auction-api/internal/lib"
)

func AuthGuard(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")

		if reqToken == "" {
			helper.WriteResponseMessage(w, "You need to be logged in to access this", http.StatusBadRequest)
			return
		}

		splitToken := strings.Split(reqToken, "Bearer")

		if len(splitToken) != 2 {
			helper.WriteResponseMessage(w, "Bearer token not in proper format", http.StatusBadRequest)
			return
		}

		reqToken = strings.TrimSpace(splitToken[1])

		isValid, sub, err := lib.VerifyToken(reqToken)

		if errors.Is(err, lib.ErrUnknownClaims) {
			helper.WriteResponseMessage(w, "Unknown token claims", http.StatusBadRequest)
			return
		}

		if errors.Is(err, lib.ErrTokenExpired) {
			helper.WriteResponseMessage(w, "Authorization token is expired", http.StatusBadRequest)
			return
		}

		if !isValid {
			helper.WriteResponseMessage(w, "Invalid token", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), helper.UserIDKey, sub)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
