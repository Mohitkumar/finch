package httpserver

import (
	"encoding/json"
	"net/http"
)

func HandleCreateFlow(w http.ResponseWriter, r *http.Request) error {
	var flow Flow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
