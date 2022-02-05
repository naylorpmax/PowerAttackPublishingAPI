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

	// temporary static data
	object := &spell.Object{
		Name:        "Borov's Golden Compass",
		Level:       "3",
		School:      "Divination",
		CastingTime: "1 action (ritual)",
		Range:       "Self",
		Components:  "V, M (a golden coin)",
		Duration:    "Instantaneous",
		Classes:     []string{"Bard", "Cleric", "Wizard"},
		Source:      "Other Homebrew Spells",
		Description: "You ask a question and carve two different answers on the golden coin, one on each side. You then flip the coin, and the side with the truest answer will face upwards.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		return err
	}
	return nil
}
