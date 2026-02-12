// auth middleware with same middleware pattern as logging and cors taking handler and returning handler

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Suthar345Piyush/invoicego/internal/domain"
	"github.com/Suthar345Piyush/invoicego/internal/util"
)

type contextKey string

const UserContextKey contextKey = "user"

// auth middleware function having the jwt secret token
// same pattern taking handler and returning handler

func Auth(jwtSecret string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		//  returning the http handler , passing func as an interface

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				util.WriteError(w, http.StatusUnauthorized, domain.ErrUnauthorized)
				return
			}

			// expected format of user auth token = "Bearer <token>"

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				util.WriteError(w, http.StatusUnauthorized, domain.ErrInvalidToken)
				return
			}

			token := parts[1]
			claims, err := util.ValidateToken(token, jwtSecret)
			if err != nil {
				util.WriteError(w, http.StatusUnauthorized, domain.ErrInvalidToken)
				return
			}

			// adding claims into context

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))

		})

	}

}

// function for getting user from passed context , it will return jwt claims , and user's context key

func GetUserFromContext(ctx context.Context) (*util.JWTClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*util.JWTClaims)

	return claims, ok
}
