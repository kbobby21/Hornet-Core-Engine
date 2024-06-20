package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"golang.org/x/exp/slog"
)

func (s *server) handleAssetAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in monitor assets handler")
		ue := r.Context().Value(factory.UserEmail)
		userEmail, ok := ue.(string)
		if !ok {
			err := errors.New("type assertion failed")
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		var postData struct {
			Monitor []factory.MonitorAssetsData `json:"monitor"`
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		// Call AddAssets method to insert the received data
		err = s.assets.AddAssets(
			postData.Monitor,
			userEmail,
		)
		if err != nil {
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		s.respond(w, &ResponseMsg{Message: "success"}, http.StatusOK)
	}
}

func (s *server) handleAssetGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handleAssetGet")
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
		assets, err := s.assets.GetAssets(page_no, userEmail)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: assets}, 200)
	}
}

func (s *server) handleUserAssetsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handleUserAssets")
		ua, err := s.assets.GetUserAssets()
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return

		}
		s.respond(w, ResponseMsg{Message: "success", Data: ua}, 200)

	}
}
