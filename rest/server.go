package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/persistence"
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
	router   *mux.Router
	pFactory persistence.PersistenceFactory
}

func NewServer(config Config) (*Server, error) {
	redisConfig := persistence.RedisConfig{
		Host:      config.Host,
		Port:      config.RedisConfig.Port,
		Namespace: config.Namespace,
	}
	cnf := persistence.Config{
		RedisConfig: redisConfig,
	}
	pFactory := new(persistence.PersistenceFactory)
	pFactory.Init(cnf, persistence.REDIS_PERSISTENCE_IMPL)
	s := &Server{
		Config:   config,
		router:   mux.NewRouter(),
		pFactory: *pFactory,
	}

	return s, nil
}

func (s *Server) Start() error {
	logger.Info("startting http server on", zap.Int("port", s.Port))
	s.router.HandleFunc("/workflow", s.HandleCreateFlow).Methods(http.MethodPost)
	s.router.HandleFunc("/workflow/{name}", s.HandleGetFlow).Methods(http.MethodGet)

	s.router.Use(loggingMiddleware)
	http.Handle("/", s.router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.router); err != nil {
		return err
	}
	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondOK(w http.ResponseWriter, message string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res, _ := json.Marshal(map[string]string{"message": message})
	w.Write(res)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
