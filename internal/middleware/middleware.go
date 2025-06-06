package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/andras-szesztai/fem_fitness_project/internal/store"
	"github.com/andras-szesztai/fem_fitness_project/internal/tokens"
	"github.com/andras-szesztai/fem_fitness_project/internal/utils"
)

type UserMiddleware struct {
	userStore store.UserStore
}

func NewUserMiddleware(userStore store.UserStore) *UserMiddleware {
	return &UserMiddleware{userStore: userStore}
}

type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		return store.AnonymousUser
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.userStore.GetUserToken(tokens.ScopeAuthentication, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid or expired token"})
			return
		}

		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid or expired token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
