package middleware

//intercept requests before they go through for validation
import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zidariu-sabin/femProject/internal/store"
	"github.com/zidariu-sabin/femProject/internal/tokens"
	"github.com/zidariu-sabin/femProject/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

// we define an extra type for the context key so that we don't get colisions on the string type in the context
type contextKey string

const UserContextKey = contextKey("user")

func NewUserMiddleware(userStore *store.PostgresUserStore) *UserMiddleware {
	return &UserMiddleware{UserStore: userStore}
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	//retrieves key and ensures if it matches type store.User
	user, ok := r.Context().Value(UserContextKey).(*store.User)

	if !ok {
		panic("missing user in request")
	}

	return user
}

// func setUserCookie(r *http.Request, user *store.User) *http.Request

// we use this function to wrap all handlers that process requests
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//whithin this anonymous function we can interject any incomign requests to the server

		//The Vary HTTP header is used to inform caches about which request headers influence the response content
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		fmt.Printf("authHeader:%v \n", authHeader)
		// Parsing auth token of type "Bearer <AUTH TOKEN"
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)

		// fmt.Printf("user:%v \n", user)
		if err != nil {
			fmt.Printf("error:%v \n", err)
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.WriteJson(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
