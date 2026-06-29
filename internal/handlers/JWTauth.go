package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	TgID int64 `json:"sub"`
	jwt.RegisteredClaims
}

func (e *Env) ValidateJWT(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jwtraw := r.Header.Get("Authorization")
		JWT := strings.TrimPrefix(jwtraw, "Bearer ")
		token, err := jwt.ParseWithClaims(JWT, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(e.JwtSecret), nil
		})
		if err != nil {
			log.Println("Error in validation JWT:", err)
			w.WriteHeader(401)
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			w.WriteHeader(401)
			return
		}

		ctx := context.WithValue(r.Context(), "TgID", claims.TgID)
		newReq := r.WithContext(ctx)
		next.ServeHTTP(w, newReq)

	})
}
