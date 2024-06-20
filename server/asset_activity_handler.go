package server

import (
	"errors"
	"net/http"
	"strconv"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"golang.org/x/exp/slog"
)

func (s *server) handleAssetActivityGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handleAssetActivityGet")
		ue := r.Context().Value(factory.UserEmail)
		userEmail, ok := ue.(string)
		if !ok {
			err := errors.New("type assertion failed")
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		page_no, err := strconv.Atoi(r.URL.Query().Get("page_no"))
		if err != nil {
			page_no = 1
		}

		activities, err := s.notifications.GetAssetActivities(page_no, userEmail)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: activities}, http.StatusOK)
	}
}
