package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/naylorpmax/homebrew-users-api/pkg/middleware"
	"github.com/naylorpmax/homebrew-users-api/pkg/patreon"
)

type (
	Welcome struct {
		OAuth2Config *oauth2.Config
	}
)

func (wel *Welcome) Handler(w http.ResponseWriter, r *http.Request) error {
	patreonClient, err := patreon.New(r, wel.OAuth2Config)
	if err != nil {
		return &middleware.Error{
			StatusCode: http.StatusForbidden,
			Message:    errors.New("unable to authenticate to Patreon").Error(),
			Details:    err.Error(),
		}
	}

	userName, err := patreonClient.AuthenticateUser()
	if err != nil {
		return &middleware.Error{
			StatusCode: http.StatusForbidden,
			Message:    errors.New("unable to authenticate user").Error(),
			Details:    err.Error(),
		}
	}

	welcomeMsg := map[string]string{
		"message": "welcome! you're logged in",
		"name":    userName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(welcomeMsg)
	if err != nil {
		return err
	}
	return nil
}
