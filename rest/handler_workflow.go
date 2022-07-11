package rest

import (
	"encoding/json"
	"net/http"

	"github.com/mohitkumar/finch/logger"
	"github.com/mohitkumar/finch/model"
	"go.uber.org/zap"
)

func (s *Server) HandleRunFlow(w http.ResponseWriter, r *http.Request) {
	var runReq model.WorkflowRunRequest
	if err := json.NewDecoder(r.Body).Decode(&runReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer r.Body.Close()
	err := s.executorService.StartFlow(runReq.Name, runReq.Data)
	if err != nil {
		logger.Error("error running workflow", zap.String("name", runReq.Name), zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "error running workflow")
		return
	}
	respondOK(w, "accepted")
}
