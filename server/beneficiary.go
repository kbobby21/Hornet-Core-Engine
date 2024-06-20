package server

import (
	"net/http"

	"golang.org/x/exp/slog"
)

func (s *server) handleBeneficiary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in Get Beneficiary")

		sender := r.URL.Query().Get("sender")

		beneficiarySummaries, err := s.beneficiary.GetBeneficiary(sender)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: beneficiarySummaries}, 200)
	}
}
