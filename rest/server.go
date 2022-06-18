package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence"
	"github.com/mohitkumar/finch/persistence/redis"
	"go.uber.org/zap"
)

type Config struct {
	RedisConfig
	Port int
}

type RedisConfig struct {
	Host      string
	Port      uint16
	Namespace string
}
type Server struct {
	Config
	router *mux.Router
	wfDao  persistence.WorkflowDao
}

func NewServer(config Config) (*Server, error) {
	cnf := redis.Config{
		Host:      config.Host,
		Port:      config.RedisConfig.Port,
		Namespace: config.Namespace,
	}
	s := &Server{
		Config: config,
		router: mux.NewRouter(),
		wfDao:  redis.NewRedisWorkflowDao(cnf),
	}

	return s, nil
}

func (s *Server) Start() error {
	logger.Info("startting http server on", zap.Int("port", s.Port))
	s.router.HandleFunc("/workflow", s.HandleCreateFlow).Methods(http.MethodPost)
	s.router.HandleFunc("/workflow/{name}", s.HandleGetFlow).Methods(http.MethodGet)
	http.Handle("/", s.router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.router); err != nil {
		return err
	}
	return nil
}
