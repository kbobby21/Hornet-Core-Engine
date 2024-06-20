package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

// IsAuthenticated checks if the given request has valid JWT token
func (s *server) IsAuthenticated(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			slog.Error("Error in getting cookie", err)
			s.respond(w, &ResponseMsg{Message: "No Cookie: " + err.Error()}, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(c.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}

			return []byte(viper.GetString("jwt_secret")), nil
		})
		if err != nil {
			slog.Error("Error in parsing cookie" + err.Error())
			s.respond(w, &ResponseMsg{Message: "Error in parsing cookie" + err.Error()}, http.StatusBadRequest)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !token.Valid || !ok {
			slog.Error("Error in parsing claims or invalid token")
			s.respond(w, &ResponseMsg{Message: "Error in parsing cookie" + err.Error()}, http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		slog.Info("Logged in user", "email", claims["email"])

		// Set the email as a key in the context
		ctx := context.WithValue(r.Context(), factory.UserEmail, claims["email"])

		// run the wrapped handler
		h(w, r.WithContext(ctx))
	}
}

func (s *server) enableCors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", viper.GetString("base_url"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, X-Api-Key, X-Requested-With , Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

func (s *server) IsAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Running in IsAdmin middleware")

		token := r.URL.Query().Get("token")
		var userEmail string

		if token != "" {
			isAdmin, email, err := s.admin.CheckAdminAndValidToken(token)
			if err != nil {
				if strings.Contains(err.Error(), "no rows in the result set") {
					s.respond(w, &ResponseMsg{Message: "User is not an admin or the token is invalid"}, http.StatusUnauthorized)
				} else {
					s.logger.Error("Error in checking admin status or token validation:", err)
					s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
				}
				return
			}

			if isAdmin {
				userEmail = email
				s.logger.Info("Logged in as admin", "email", userEmail)
			} else {
				s.respond(w, &ResponseMsg{Message: "Admin access required"}, http.StatusUnauthorized)
				return
			}
		} else {
			c, err := r.Cookie("token")
			if err != nil {
				slog.Error("Error in getting cookie", err)
				s.respond(w, &ResponseMsg{Message: "No Cookie: " + err.Error()}, http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(c.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
				}

				return []byte(viper.GetString("jwt_secret")), nil
			})
			if err != nil {
				slog.Error("Error in parsing cookie" + err.Error())
				s.respond(w, &ResponseMsg{Message: "Error in parsing cookie" + err.Error()}, http.StatusBadRequest)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !token.Valid || !ok {
				slog.Error("Error in parsing claims or invalid token")
				s.respond(w, &ResponseMsg{Message: "Error in parsing cookie" + err.Error()}, http.StatusBadRequest)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			userEmail = claims["email"].(string)

			isAdmin, err := s.admin.CheckAdminAndValidTokenByEmail(userEmail)
			if err != nil {
				s.logger.Error("Error in checking admin status by email:", err)
				s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
				return
			}

			if isAdmin {
				s.logger.Info("Logged in as admin")
			} else {
				s.respond(w, &ResponseMsg{Message: "Admin access required"}, http.StatusUnauthorized)
				return
			}
		}
		ctx := context.WithValue(r.Context(), factory.UserEmail, userEmail)

		handler(w, r.WithContext(ctx))
	}
}
