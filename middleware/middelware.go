package middleware

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// authHeader := r.Header.Get("Authorization")
		// if authHeader == "" {
		// 	http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		// 	return
		// }

		// tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		// if tokenStr == authHeader {
		// 	http.Error(w, "Bearer token missing", http.StatusUnauthorized)
		// 	return
		// }
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized - no session found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid cookie", http.StatusBadRequest)
			return
		}

		sessionToken := cookie.Value

		secretKey := []byte(os.Getenv("Secret_Key"))
		token, err := jwt.Parse(sessionToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
