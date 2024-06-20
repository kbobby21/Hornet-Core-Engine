package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"golang.org/x/exp/slog"
)

func (s *server) handleNotificationsAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in notificationsHandler")

		var notifications []factory.NotificationsData
		err := json.NewDecoder(r.Body).Decode(&notifications)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		err = s.notifications.AddNotification(notifications)
		if err != nil {
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}

func (s *server) handleNotificationsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handleNotificationsGet")
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

		notifications, err := s.notifications.GetNotifications(page_no, userEmail)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: notifications}, 200)

	}
}

func (s *server) handleNotificationsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in Delete Notifications")

		ue := r.Context().Value(factory.UserEmail)
		userEmail, ok := ue.(string)
		if !ok {
			err := errors.New("User email not found")
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusUnauthorized)
			return
		}

		err := s.notifications.DeleteNotifications(userEmail)
		if err != nil {
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}
