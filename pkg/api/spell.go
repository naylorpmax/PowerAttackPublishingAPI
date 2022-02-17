package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/naylorpmax/homebrew-users-api/pkg/middleware"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
)

type (
	SpellRequest struct {
		Name *string
	}

	SpellLookup struct {
		SpellService *spell.Service
	}
)

func (s *SpellLookup) Handler(w http.ResponseWriter, r *http.Request) error {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		return &middleware.Error{
			Message:    "unsupported media type",
			Details:    fmt.Sprintf("expected 'application/json', got '%v'", contentType),
			StatusCode: http.StatusBadRequest,
		}
	}

	spellReq := &SpellRequest{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(spellReq); err != nil {
		return &middleware.Error{
			Message:    "unable to unmarshal request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	if spellReq.Name == nil || *spellReq.Name == "" {
		return &middleware.Error{
			Message:    "missing required body element 'name'",
			StatusCode: http.StatusBadRequest,
		}
	}

	spells, err := s.SpellService.Lookup(r.Context(), *spellReq.Name)
	if err != nil {
		return &middleware.Error{
			Message:    "unable to lookup spell",
			StatusCode: http.StatusInternalServerError,
			Details:    err.Error(),
		}
	}

	results := map[string][]*spell.Object{
		"spells": spells,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		return err
	}
	return nil
}
