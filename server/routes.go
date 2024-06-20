package server

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/hornetdefiant/core-engine/pkg/admin"
	"bitbucket.org/hornetdefiant/core-engine/pkg/analyze"
	"bitbucket.org/hornetdefiant/core-engine/pkg/assets"
	"bitbucket.org/hornetdefiant/core-engine/pkg/beneficiary"
	"bitbucket.org/hornetdefiant/core-engine/pkg/customers"
	"bitbucket.org/hornetdefiant/core-engine/pkg/darkweb"
	"bitbucket.org/hornetdefiant/core-engine/pkg/exchange"
	"bitbucket.org/hornetdefiant/core-engine/pkg/mail"
	"bitbucket.org/hornetdefiant/core-engine/pkg/notifications"
	"bitbucket.org/hornetdefiant/core-engine/pkg/sessmanager"
	"bitbucket.org/hornetdefiant/core-engine/pkg/tags"
	"bitbucket.org/hornetdefiant/core-engine/pkg/token"
	"bitbucket.org/hornetdefiant/core-engine/pkg/transactions"
	"bitbucket.org/hornetdefiant/core-engine/pkg/users"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

type server struct {
	router        *mux.Router
	logger        *slog.Logger
	customer      customers.Repository
	sessmanager   sessmanager.Repository
	user          users.Repository
	darkweb       darkweb.Repository
	exchange      exchange.Repository
	mail          mail.Repository
	txs           transactions.Repo
	tags          tags.Repository
	notifications notifications.Repository
	assets        assets.Repository
	token         token.Repository
	beneficiary   beneficiary.Repository
	analyze       analyze.Repository
	admin         admin.Repository
}

type ResponseMsg struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Define all the routes here
func (s *server) routes() {
	// Add details of a darkwebsite
	s.router.HandleFunc(
		"/darkwebsite",
		s.enableCors(s.IsAdmin(s.addDarkWebMeta())),
	).Methods(http.MethodPost, http.MethodOptions)

	// Get details of a darkwebsite
	s.router.HandleFunc(
		"/darkwebsite",
		s.enableCors(s.IsAuthenticated(s.getDarkWebMeta())),
	).Methods(http.MethodGet, http.MethodOptions)

	// Add details of the interested customers
	s.router.HandleFunc(
		"/customer",
		s.enableCors(s.handleCustomerAdd()),
	).Methods(http.MethodPost, http.MethodOptions)

	// Add metadata of exchange
	s.router.HandleFunc(
		"/exchangemeta",
		s.enableCors(s.IsAdmin(s.addExchangeMetaData())),
	).Methods(http.MethodPost, http.MethodOptions)

	// fetch the metadata of exchange
	s.router.HandleFunc(
		"/exchangemeta",
		s.enableCors(s.IsAuthenticated(s.getExchangeMetaData())),
	).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc(
		"/signup",
		s.enableCors(s.handleSignUp()),
	).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc(
		"/verify",
		s.enableCors(s.handleVerify()),
	).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc(
		"/signin",
		s.enableCors(s.handleSignIn()),
	).Methods(http.MethodPost, http.MethodOptions)

	//Get method for users
	s.router.HandleFunc(
		"/users",
		s.enableCors(s.IsAdmin(s.handleGetAllUsers())),
	).Methods(http.MethodGet, http.MethodOptions)

	//Get method for users
	s.router.HandleFunc(
		"/user",
		s.enableCors(s.IsAdmin(s.handleGetUser())),
	).Methods(http.MethodGet, http.MethodOptions)

	// Add transaction details to core-engine
	s.router.HandleFunc(
		"/transactions",
		s.enableCors(s.IsAdmin(s.handleTxAdd())),
	).Methods(http.MethodPost, http.MethodOptions)

	// Receive tx details from txmonitor
	s.router.HandleFunc(
		"/transactions",
		s.enableCors(s.IsAuthenticated(s.handleTxRetrieve())),
	).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc(
		"/wallet_transactions",
		s.enableCors(s.IsAuthenticated(s.handleWalletTxs())),
	).Methods(http.MethodGet, http.MethodOptions)

	// Fetch the balance of a wallet
	s.router.HandleFunc(
		"/wallet_balance",
		s.enableCors(s.IsAuthenticated(s.handleWalletBalance())),
	).Methods(http.MethodGet, http.MethodOptions)

	//Retrieve tags data
	s.router.HandleFunc(
		"/tags_graph",
		s.enableCors(s.IsAuthenticated(s.handleTagsRetrieve())),
	).Methods(http.MethodGet, http.MethodOptions)

	//Add data for notifications
	s.router.HandleFunc(
		"/notifications",
		s.enableCors(s.IsAdmin(s.handleNotificationsAdd())),
	).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc(
		"/notifications",
		s.enableCors(s.IsAuthenticated(s.handleNotificationsGet())),
	).Methods(http.MethodGet, http.MethodOptions)

	//delete notifications
	s.router.HandleFunc(
		"/notifications",
		s.enableCors(s.IsAuthenticated(s.handleNotificationsDelete())),
	).Methods(http.MethodDelete, http.MethodOptions)

	// Add data of the assests to be monitored
	s.router.HandleFunc("/assets",
		s.enableCors(s.IsAuthenticated(s.handleAssetAdd())),
	).Methods(http.MethodPost, http.MethodOptions)

	// Get data for assets
	s.router.HandleFunc("/assets",
		s.enableCors(s.IsAuthenticated(s.handleAssetGet())),
	).Methods(http.MethodGet, http.MethodOptions)

	//get asset activity
	s.router.HandleFunc("/asset_activity",
		s.enableCors(s.IsAuthenticated(s.handleAssetActivityGet())),
	).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc("/user_assets",
		s.IsAdmin(s.handleUserAssetsGet()),
	).Methods(http.MethodGet)

	// Store token
	s.router.HandleFunc("/token_store",
		s.enableCors(s.IsAuthenticated(s.handleTokenStore())),
	).Methods(http.MethodPost, http.MethodOptions)

	//get beneficiary
	s.router.HandleFunc("/beneficiary",
		s.enableCors(s.IsAuthenticated(s.handleBeneficiary())),
	).Methods(http.MethodGet, http.MethodOptions)

	//get analyze wallet transaction
	s.router.HandleFunc("/analyze_transaction",
		s.enableCors(s.IsAuthenticated(s.handleAnalyzeTransaction())),
	).Methods(http.MethodGet, http.MethodOptions)

	// serve the react build dir
	s.router.HandleFunc(
		"/static/",
		s.serveReactStaticFiles(),
	)
	s.router.HandleFunc(
		"/",
		s.serveReactIndex(),
	)

}

// Respond to all API calls
func (s *server) respond(
	w http.ResponseWriter,
	data interface{},
	status int,
) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.logger.Error("Error in encoding the response", "error", err)
			return
		}
	}
}
