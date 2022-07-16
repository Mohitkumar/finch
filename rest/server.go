package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mohitkumar/finch/container"
	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/service"
	"go.uber.org/zap"
)

type Server struct {
	HttpPort        int
	router          *mux.Router
	container       *container.DIContiner
	executorService *service.WorkflowExecutionService
}

func NewServer(httpPort int, container *container.DIContiner) (*Server, error) {

	s := &Server{
		HttpPort:        httpPort,
		router:          mux.NewRouter(),
		container:       container,
		executorService: service.NewWorkflowExecutionService(container),
	}

	return s, nil
}

func (s *Server) Start() error {
	logger.Info("startting http server on", zap.Int("port", s.HttpPort))
	s.router.HandleFunc("/workflow", s.HandleCreateFlow).Methods(http.MethodPost)
	s.router.HandleFunc("/workflow/{name}", s.HandleGetFlow).Methods(http.MethodGet)
	s.router.HandleFunc("/flow/execute", s.HandleRunFlow).Methods(http.MethodPost)
	s.router.Use(loggingMiddleware)
	http.Handle("/", s.router)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.HttpPort), s.router); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
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
