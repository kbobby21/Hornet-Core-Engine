package server

import (
	"net/http"

	"golang.org/x/exp/slog"
)

func (s *server) handleAnalyzeTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in Analyze wallet transaction")

		traddress := r.URL.Query().Get("traddress")
		coaddress := r.URL.Query().Get("coaddress")

		transactionInfo, err := s.analyze.GetTransactionInfo(traddress, coaddress)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}

		s.respond(w, ResponseMsg{Message: "success", Data: transactionInfo}, 200)
	}
}
