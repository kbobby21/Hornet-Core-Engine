package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"golang.org/x/exp/slog"
)

func (s *server) handleTokenStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in Token Store handler")
		ue := r.Context().Value(factory.UserEmail)
		userEmail, ok := ue.(string)
		if !ok {
			err := errors.New("type assertion failed")
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		var postData struct {
			Token factory.Tokens `json:"token"`
		}

		err := json.NewDecoder(r.Body).Decode(&postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}

		err = s.token.AddToken(
			postData.Token,
			userEmail,
		)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success"}, 200)
	}
}
