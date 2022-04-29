package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Config struct {
	Port int
}
type Server struct {
	Config
	router *mux.Router
	logger *zap.Logger
}

func NewServer(config Config) (*Server, error) {
	s := &Server{
		Config: config,
		router: mux.NewRouter(),
		logger: zap.L().Named("httpserver"),
	}

	return s, nil
}

func (s *Server) Start() error {
	s.logger.Info("startting http server on", zap.Int("port", s.Port))
	s.router.HandleFunc("/flow", s.HandleCreateFlow).Methods(http.MethodPost)
	s.router.HandleFunc("/flow/{name}", s.HandleGetFlow).Methods(http.MethodGet)
	http.Handle("/", s.router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.router); err != nil {
		return err
	}
	return nil
}
