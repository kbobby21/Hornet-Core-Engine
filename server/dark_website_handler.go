package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (s *server) addDarkWebMeta() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.logger.Info("Running in addDarkWebMeta")
		var postData factory.DarkWebSite
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		err = s.darkweb.AddWebSiteMeta(postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}

func (s *server) getDarkWebMeta() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in getDarkWebMeta")
		pn, err := strconv.Atoi(r.URL.Query().Get("page_no"))
		if err != nil {
			// default page number is 1
			pn = 1
		}
		dw, err := s.darkweb.GetWebsiteMeta(pn)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
		}
		s.respond(w, &ResponseMsg{Message: "success", Data: dw}, 200)
	}
}
