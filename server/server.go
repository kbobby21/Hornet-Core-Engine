package server

import (
	"log"
	"net/http"
	"os"

	"bitbucket.org/hornetdefiant/core-engine/db/postgres"
	"bitbucket.org/hornetdefiant/core-engine/db/redis"
	"bitbucket.org/hornetdefiant/core-engine/pkg/mail"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func Run(envType *string) {
	// Setup the configuration management
	viper.SetConfigType("json") // type of the config file

	var level slog.Level
	if *envType == "dev" {
		viper.SetConfigName("config-dev") // name of the config file
		level = slog.LevelDebug

	} else if *envType == "prod" {
		viper.SetConfigName("config-prod") // name of the config file
		level = slog.LevelInfo

	} else {
		viper.SetConfigName("config-staging") // name of the config file
		level = slog.LevelDebug
	}

	viper.AddConfigPath("$HOME/.hornet") // location of local/dev config file

	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Error in reading config file", err)
		return
	}

	// add options to the logger
	ho := slog.HandlerOptions{
		AddSource: true, // prints the line in code
		Level:     level,
	}

	postgres := postgres.NewPostgres()
	redis := redis.NewRedis()
	// &slog.HandlerOptions{AddSource: true, ReplaceAttr: replace}
	serve := &server{
		router:      mux.NewRouter(),
		logger:      slog.New(slog.NewJSONHandler(os.Stdout, &ho)),
		customer:    postgres,
		sessmanager: redis,
		user:        postgres,
		darkweb:     postgres,
		exchange:    postgres,
		mail: mail.NewZohoMail(
			viper.GetString("zoho_mail_id"),
			viper.GetString("zoho_mail_password"),
		),
		txs:           postgres,
		tags:          postgres,
		notifications: postgres,
		assets:        postgres,
		token:         postgres,
		beneficiary:   postgres,
		analyze:       postgres,
		admin:         postgres,
	}

	// setup the routes
	serve.routes()

	slog.Info("Serving react webapp on", "port", viper.GetString("reactServerPort"))
	slog.Info("Serving API server on", "port", viper.GetString("apiServerPort"))

	// Serve react static files
	http.HandleFunc("/static/", serve.serveReactStaticFiles())

	// Serve index.html for all other routes
	http.HandleFunc("/", serve.serveReactIndex())

	if viper.GetString("reactServerPort") == "443" {
		go func() {
			log.Fatal(http.ListenAndServeTLS(
				":"+viper.GetString("reactServerPort"),
				viper.GetString("fullChainPath"),
				viper.GetString("privKeyPath"),
				nil,
			))
		}()

		log.Fatal(
			http.ListenAndServeTLS(
				":"+viper.GetString("apiServerPort"),
				viper.GetString("fullChainPath"),
				viper.GetString("privKeyPath"),
				serve,
			).Error(),
		)

	} else {
		go func() {
			log.Fatal(http.ListenAndServe(
				":"+viper.GetString("reactServerPort"),
				nil,
			))
		}()

		log.Fatal(
			http.ListenAndServe(
				":"+viper.GetString("apiServerPort"),
				serve,
			).Error(),
		)

	}
}
