package middleware

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"

	"payments_service/ctxutils"
	"payments_service/services"
)

type AuthMiddleware struct {
	service *services.TokenStruct
}

func NewAuthMiddleware(service *services.TokenStruct) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
	}
}

func (a AuthMiddleware) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
		log.Debug().Msgf("tokenString: %s", tokenString)

		//check
		claims, err := a.service.ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		//input UserID in context
		ctx := ctxutils.WithUserID(r.Context(), claims.User_id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
