package api

import (
	"net/http"

	"github.com/dorik33/medods_test/internal/auth"
	"github.com/dorik33/medods_test/internal/config"
	"github.com/dorik33/medods_test/internal/handlers"
	"github.com/dorik33/medods_test/internal/middleware"
	"github.com/dorik33/medods_test/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type API struct {
	config   *config.Config
	router   *mux.Router
	logger   *logrus.Logger
	store    *store.Store
	handlers *handlers.Handlers
	auth     *auth.Service
}

func New(cfg *config.Config) *API {
	api := &API{
		config: cfg,
		router: mux.NewRouter(),
		logger: logrus.New(),
	}
	return api
}

func (api *API) Start() error {
	if err := api.configureLogger(); err != nil {
		return err
	}

	dbStore, err := store.NewConnection(api.config, api.logger)
	if err != nil {
		return err
	}
	api.store = dbStore

	api.logger.Debug("Successful connection to database")

	defer api.store.Close()

	api.auth = auth.NewService(
		api.store,
		api.config,
	)

	api.handlers = handlers.NewHandlers(api.logger, api.store, api.auth)
	api.configureRouter()
	server := &http.Server{
		Handler:      api.router,
		Addr:         api.config.Addr,
		WriteTimeout: api.config.WriteTimeout,
	}

	api.logger.Debug("Server is running with addr: ", api.config.Addr)

	return server.ListenAndServe()
}

func (api *API) configureLogger() error {
	api.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	level, err := logrus.ParseLevel(api.config.LogLevel)
	if err != nil {
		return err
	}
	api.logger.SetLevel(level)

	return nil
}

func (api *API) configureRouter() {
	api.router.Use(middleware.JSONContentTypeMiddleware)
	api.router.Use(middleware.LoggingMiddleware(api.logger))
	api.router.HandleFunc("/auth/token", api.handlers.GenerateTokensHandler).Methods("GET")
	api.router.HandleFunc("/auth/token/refresh", api.handlers.RefreshTokenHandler).Methods("POST")
}
