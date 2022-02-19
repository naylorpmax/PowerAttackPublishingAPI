package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	pat "github.com/naylorpmax/homebrew-users-api/pkg/client/patreon"
	"github.com/naylorpmax/homebrew-users-api/pkg/middleware/apierror"
	"github.com/naylorpmax/homebrew-users-api/pkg/patreon"
)

type (
	Welcome struct {
		OAuth2Config *oauth2.Config
	}
)

func (wel *Welcome) Handler(w http.ResponseWriter, r *http.Request) error {
	code := r.FormValue("code")
	if code == "" {
		return &apierror.Error{
			StatusCode: http.StatusForbidden,
			Message:    "unable to authenticate to Patreon",
			Details:    "redirect request does not contain OAuth2 code",
		}
	}

	client, err := pat.New(r.Context(), code, wel.OAuth2Config)
	if err != nil {
		return &apierror.Error{
			StatusCode: http.StatusForbidden,
			Message:    "unable to authenticate to Patreon",
			Details:    err.Error(),
		}
	}

	patreonClient, err := patreon.New(client)
	if err != nil {
		return &apierror.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    errors.New("unable to create Patreon client").Error(),
			Details:    err.Error(),
		}
	}

	userName, err := patreonClient.AuthenticateUser()
	if err != nil {
		return &apierror.Error{
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
