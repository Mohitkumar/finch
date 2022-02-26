package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/mohitkumar/finch/api/v1"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func (s *Server) HandleCreateFlow(w http.ResponseWriter, r *http.Request) {
	var flow Flow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	ctx := context.Background()
	data := make(map[string]*structpb.Value)
	for k, v := range flow.Data {
		if val, err := structpb.NewValue(v); err == nil {
			data[k] = val
		}
	}
	var actions []*api.Action
	for _, action := range flow.Actions {
		inputParams := action.InputParameters
		inputParamsPb := make(map[string]*structpb.Value)
		for k, v := range inputParams {
			if val, err := structpb.NewValue(v); err == nil {
				inputParamsPb[k] = val
			}
		}
		actions = append(actions, &api.Action{
			Id:              action.Id,
			Name:            action.Name,
			InputParameters: inputParamsPb,
		})
	}
	res, err := s.getCoordClient().CreateFlow(ctx, &api.FlowCreateRequest{
		Flow: &api.Flow{
			Name:          flow.Name,
			Data:          data,
			StartActionId: flow.StartActionId,
			Actions:       actions,
		},
	})
	if err != nil {
		log.Fatal(err)
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
	ctx := context.Background()
	res, err := s.getCoordClient().GetFlow(ctx, &api.FlowGetRequest{
		Name: flowName,
	})
	if err != nil {
		log.Fatal(err)
	}
	flow := res.Flow
	data := make(map[string]interface{})
	for k, v := range flow.Data {
		if val := v.AsInterface(); err == nil {
			data[k] = val
		}
	}
	var actions []Action
	for _, action := range flow.Actions {
		inputParams := action.InputParameters
		inputParamsPb := make(map[string]interface{})
		for k, v := range inputParams {
			if val := v.AsInterface(); err == nil {
				inputParamsPb[k] = val
			}
		}
		actions = append(actions, Action{
			Id:              action.Id,
			Name:            action.Name,
			InputParameters: inputParamsPb,
		})
	}
	response := &Flow{
		Name:          flow.Name,
		StartActionId: flow.StartActionId,
		Data:          data,
		Actions:       actions,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
