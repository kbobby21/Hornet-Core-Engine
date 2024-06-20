package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (s *server) addExchangeMetaData() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in exmonitorHandler")
		var postData factory.ExchangeMetaDataInsert
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		err = s.exchange.AddExchangeMetaData(&postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}

func (s *server) getExchangeMetaData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in exmonitorHandler")
		pn, err := strconv.Atoi(r.URL.Query().Get("page_no"))
		if err != nil {
			pn = 1
		}
		exchangeData, err := s.exchange.GetExchangeMetaData(pn)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, map[string]interface{}{
			"message": "success",
			"data":    exchangeData,
		}, 200)
	}
}
