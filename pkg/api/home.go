package api

import (
	"encoding/json"
	"net/http"
)

type (
	Home struct{
	}
)

func (h *Home) Handler(w http.ResponseWriter, r *http.Request) error {
	// temporary static data
	respBody := map[string]string{
		"message": "you're home!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(respBody)
	if err != nil {
		return err
	}
	return nil
}
