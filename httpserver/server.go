package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/mohitkumar/finch/api/v1"
	_ "github.com/mohitkumar/finch/loadbalance"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Port         int
	CoordRpcPort int
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

func (s *Server) getCoordClient() api.CoordinatorClient {
	conn, err := grpc.Dial(fmt.Sprintf("coordinator:///127.0.0.1:%d", s.CoordRpcPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return api.NewCoordinatorClient(conn)
}
