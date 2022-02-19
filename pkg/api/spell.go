package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/naylorpmax/homebrew-users-api/pkg/middleware/apierror"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
)

type (
	SpellRequest struct {
		Name  *string `json:"name"`
		Level *string `json:"level"`
	}

	SpellLookup struct {
		SpellService *spell.Service
	}
)

func (s *SpellLookup) Handler(w http.ResponseWriter, r *http.Request) error {
	errCh := make(chan error)

	go func() {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			errCh <- &apierror.Error{
				Message:    "unsupported media type",
				Details:    fmt.Sprintf("expected 'application/json', got '%v'", contentType),
				StatusCode: http.StatusBadRequest,
			}
			return
		}

		spellReq := &SpellRequest{}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(spellReq); err != nil {
			errCh <- &apierror.Error{
				Message:    "unable to unmarshal request",
				Details:    err.Error(),
				StatusCode: http.StatusBadRequest,
			}
			return
		}

		if spellReq.Name == nil && spellReq.Level == nil {
			errCh <- &apierror.Error{
				Message:    "request has no non-empty body properties",
				StatusCode: http.StatusBadRequest,
			}
			return
		}

		spells, err := s.SpellService.Lookup(r.Context(), spellReq.Name, spellReq.Level)
		if err != nil {
			errCh <- &apierror.Error{
				Message:    "unable to lookup spell",
				StatusCode: http.StatusInternalServerError,
				Details:    err.Error(),
			}
			return
		}

		results := map[string][]*spell.Object{
			"spells": spells,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}
