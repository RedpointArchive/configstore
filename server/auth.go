package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type c string

const contextSubjectKey c = "sub"

func authGetUser(w http.ResponseWriter, r *http.Request, signingSecret []byte, requiredScopes []string, iss string, aud string) (string, error) {
	if authorizationHeader := r.Header.Get("Authorization"); authorizationHeader == "" {
		return "", fmt.Errorf("missing 'Authorization' header")
	} else {
		token, err := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return signingSecret, nil
		})

		if err != nil {
			return "", fmt.Errorf("failed to parse JWT: %v", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return "", fmt.Errorf("invalid token claims")
		}

		scopes := strings.Split(claims["scope"].(string), " ")
		for _, rscope := range requiredScopes {
			found := false
			for _, scope := range scopes {
				if scope == rscope {
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("missing scope: %s", rscope)
			}
		}

		if claims["iss"].(string) != iss ||
			claims["aud"].(string) != aud {
			return "", fmt.Errorf("invalid iss or aud field")
		}

		return claims["sub"].(string), nil
	}
}

func authMiddleware(handler http.Handler, signingSecret []byte, requiredScopes []string, iss string, aud string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub, err := authGetUser(w, r, signingSecret, requiredScopes, iss, aud)

		if err != nil {
			// CLI user
			fmt.Printf("authentication failed: %v\n", err)
			http.Error(w, "forbidden", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, contextSubjectKey, sub)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
