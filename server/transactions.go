package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"golang.org/x/exp/slog"
)

func (s *server) handleTxAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in handleTxAdd")

		var t []factory.Tx
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		err = s.txs.AddTransactions(t)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return

		}

		s.respond(w, ResponseMsg{Message: "success"}, 200)
	}
}

func (s *server) handleTxRetrieve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handleTxRetrieve")
		pn, err := strconv.Atoi(r.URL.Query().Get("page_no"))
		if err != nil {
			// default page number is 1
			pn = 1
		}
		sender := r.URL.Query().Get("sender")
		receiver := r.URL.Query().Get("receiver")
		txs, err := s.txs.GetTransactions(pn, sender, receiver)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: txs}, 200)
	}
}

func (s *server) handleWalletTxs() http.HandlerFunc {
	type TxWithDirection struct {
		factory.Tx
		Direction string `json:"direction"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handleTxRetrieve")
		pn, err := strconv.Atoi(r.URL.Query().Get("page_no"))
		if err != nil {
			// default page number is 1
			pn = 1
		}
		allTxs := make([]TxWithDirection, 0)
		wallet := r.URL.Query().Get("wallet")
		// get all the SENT transactions
		sentTxs, err := s.txs.GetTransactions(pn, wallet, "all")
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}

		for _, tx := range sentTxs {
			var t TxWithDirection
			t.Amount = tx.Amount
			t.BlockNum = tx.BlockNum
			t.Direction = "sent to"
			t.Receiver = tx.Receiver
			t.Sender = tx.Sender
			t.Timestamp = tx.Timestamp
			allTxs = append(allTxs, t)
		}

		// get all the RECEIVED transactions
		recvdTxs, err := s.txs.GetTransactions(pn, "all", wallet)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		for _, tx := range recvdTxs {
			var t TxWithDirection
			t.Amount = tx.Amount
			t.BlockNum = tx.BlockNum
			t.Direction = "received from"
			t.Receiver = tx.Receiver
			t.Sender = tx.Sender
			t.Timestamp = tx.Timestamp
			allTxs = append(allTxs, t)
		}
		s.respond(w, ResponseMsg{Message: "success", Data: allTxs}, 200)

	}
}
func (s *server) handleWalletBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance float64
		slog.Info("Running in handleTxRetrieve")
		pn, err := strconv.Atoi(r.URL.Query().Get("page_no"))
		if err != nil {
			// default page number is 1
			pn = 1
		}
		wallet := r.URL.Query().Get("wallet")

		// get all the RECEIVED transactions
		recvdTxs, err := s.txs.GetTransactions(pn, "all", wallet)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		for _, tx := range recvdTxs {
			balance += tx.Amount
		}

		sentTxs, err := s.txs.GetTransactions(pn, wallet, "all")
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}

		for _, tx := range sentTxs {
			balance -= tx.Amount
		}
		bj := make(map[string]float64)
		bj["balance"] = balance
		s.respond(w, ResponseMsg{Message: "success", Data: bj}, 200)

	}
}
