package webHTTP

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func JWTAuth(
	original func(
		w http.ResponseWriter, 
		r *http.Request)) func(
			w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warn().Msg("received unauthorized request")
			return
		}

		authHeaderParts := strings.Split(authHeader[0]," ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warn().Msg("auth header failed to parse")
			return			
		}

		if validateToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warn().Msg("could not validate token")
			return		
		}
	}
}

func validateToken(accessToken string) bool {
	signingkey := []byte("OnceUponATime")
	
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token invalid")
		}
		
		return signingkey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}