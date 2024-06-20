package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Add new customer details from landing page
func (s *server) handleCustomerAdd() http.HandlerFunc {
	type customer struct {
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in handleCustomerAdd")
		var c customer
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil || len(c.Email) == 0 {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		// verify this user
		if err := s.mail.SendMail(
			c.Email,
			"Welcome to Hornet",
			"Thank you for subscribing to Hornet!",
		); err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		err = s.customer.CustomersAdd(c.Email)
		if err != nil {
			if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				s.respond(
					w,
					&ResponseMsg{Message: "Customer already added"},
					403,
				)
				return
			}
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 403)
			return
		}

		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}
