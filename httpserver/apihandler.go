package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	api "github.com/mohitkumar/finch/api/v1"
	"google.golang.org/grpc"
)

func (s *Server) HandleCreateFlow(w http.ResponseWriter, r *http.Request) {
	var flow Flow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	conn, err := grpc.Dial(fmt.Sprintf(":%d", s.CoordRpcPort), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewCoordinatorClient(conn)
	ctx := context.Background()
	res, err := client.CreateFlow(ctx, &api.FlowCreateRequest{
		Flow: &api.Flow{
			Name: flow.Name,
			Data: flow.Data,
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
