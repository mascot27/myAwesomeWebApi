package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}

		// get token's claims
		isSingleUse := claims.(jwt.MapClaims)["isSingleUse"].(string)
		userUuid := claims.(jwt.MapClaims)["userUuid"].(string)
		tokenUuid := claims.(jwt.MapClaims)["tokenUuid"].(string)

		// check if can be used
		// TODO: implement me


		// transmit claims after verification of single use
		r.Header.Set("isSingleUse", isSingleUse)
		r.Header.Set("userUuid", userUuid)
		r.Header.Set("tokenUuid", tokenUuid)

		next.ServeHTTP(w, r)
	})
}
