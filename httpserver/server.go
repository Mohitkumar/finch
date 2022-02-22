package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port         int
	CoordRpcPort int
}
type Server struct {
	Config
	router *mux.Router
}

func NewServer(config Config) (*Server, error) {
	return &Server{
		router: mux.NewRouter(),
	}, nil
}

func (s *Server) Start() error {
	s.router.HandleFunc("/flow/create", s.HandleCreateFlow).Methods(http.MethodPost)
	http.Handle("/", s.router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.router); err != nil {
		return err
	}
	return nil
}
