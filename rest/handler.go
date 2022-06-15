package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	api "github.com/mohitkumar/finch/api/v1"
	"github.com/mohitkumar/finch/model"
)

func (s *Server) HandleCreateFlow(w http.ResponseWriter, r *http.Request) {
	var flow model.Workflow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if res.Status == api.FlowCreateResponse_SUCCESS {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Server) HandleGetFlow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flowName, ok := vars["name"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
