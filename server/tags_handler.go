package server

import "net/http"

func (s *server) handleTagsRetrieve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in getTags")

		tags, err := s.tags.GetTags()
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		s.respond(w, &ResponseMsg{Message: "success", Data: tags}, http.StatusOK)
	}
}
