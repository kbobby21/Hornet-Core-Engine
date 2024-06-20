package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"bitbucket.org/hornetdefiant/core-engine/pkg/users"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func (s *server) handleSignUp() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in signUp Handler")

		var u factory.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 400)
			return
		}

		u.HashedPassword, err = users.HashPassword(u.Password)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		err = s.user.AddUser(u)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		// store the email and its hash for verification purpose temporarily
		hash, err := s.sessmanager.StoreEmailHash(u.Email)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.logger.Debug("Verify at", "URL", "localhost:8000/verify?ehash="+hash)
		// verify this user
		if err := s.mail.SendMail(
			u.Email,
			"Verification for Hornet",
			"Please verify your account on Hornet by clicking on this link-  "+viper.GetString("base_url")+":8000/verify?ehash="+hash,
		); err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}

func (s *server) handleVerify() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in verify Handler")

		// get the email hash
		eh := r.URL.Query().Get("ehash")

		email, err := s.sessmanager.GetEmailFromHash(eh)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		// mark this email as verified
		if err := s.user.VerifyEmail(email); err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		http.Redirect(
			w,
			r,
			viper.GetString("base_url")+"/verify",
			http.StatusSeeOther,
		)
	}
}

func (s *server) handleSignIn() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Running in signIn Handler")
		var postData factory.User
		err := json.NewDecoder(r.Body).Decode(&postData)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}
		pass, verified, err := s.user.LoginUser(postData)
		if err != nil {
			if err == sql.ErrNoRows {
				s.logger.Error("No user found, please sign up")
				s.respond(
					w,
					&ResponseMsg{Message: "No user found, please sign up"},
					http.StatusBadRequest,
				)
				return
			}
			s.logger.Error(err.Error())
			s.respond(
				w,
				&ResponseMsg{Message: err.Error()},
				500,
			)
			return
		}

		if !users.VerifyPassword(pass, postData.Password) || !verified {
			s.logger.Info("verifiaction status", "true", verified)
			s.respond(w, &ResponseMsg{Message: "Password did not match or not verified user "}, 401)
			return
		}

		token, err := users.CreateToken(postData.Email)
		if err != nil {
			s.logger.Error("Error in creating token", err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, 500)
			return
		}

		// set the cookie on client side with same expiration time
		http.SetCookie(
			w,
			&http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(10 * time.Minute),
			},
		)
		s.respond(w, &ResponseMsg{Message: "success"}, 200)
	}
}

func (s *server) handleGetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in fetch users")

		users, err := s.user.GetAllUsersEmail()
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: users}, 200)
	}
}

func (s *server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Running in handle get user")

		ue := r.Context().Value(factory.UserEmail)
		userEmail, ok := ue.(string)
		if !ok {
			err := errors.New("type assertion failed")
			s.logger.Error(err.Error())
			s.respond(w, &ResponseMsg{Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		user, err := s.user.GetUser(userEmail)
		if err != nil {
			s.logger.Error(err.Error())
			s.respond(w, ResponseMsg{Message: err.Error()}, 500)
			return
		}
		s.respond(w, ResponseMsg{Message: "success", Data: user}, 200)
	}
}
