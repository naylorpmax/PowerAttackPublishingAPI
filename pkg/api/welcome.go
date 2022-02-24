package api

import (
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/naylorpmax/gopatreon"
	"github.com/naylorpmax/homebrew-users-api/pkg/middleware/apierror"
)

type (
	Welcome struct {
		OAuth2Config *oauth2.Config
	}
)

func (wel *Welcome) Handler(w http.ResponseWriter, r *http.Request) error {
	errCh := make(chan error)

	go func() {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			errCh <- nil
			return
		}

		code := r.FormValue("code")
		if code == "" {
			errCh <- &apierror.Error{
				StatusCode: http.StatusForbidden,
				Message:    "unable to authenticate to Patreon",
				Details:    "redirect request does not contain OAuth2 code",
			}
			return
		}

		token, err := wel.OAuth2Config.Exchange(r.Context(), code)
		if err != nil {
			errCh <- &apierror.Error{
				StatusCode: http.StatusForbidden,
				Message:    "unable to authenticate to Patreon",
				Details:    err.Error(),
			}
			return
		}

		oauth2Client := wel.OAuth2Config.Client(r.Context(), token)

		client, err := gopatreon.New(oauth2Client)
		if err != nil {
			errCh <- &apierror.Error{
				StatusCode: http.StatusForbidden,
				Message:    "unable to authenticate to Patreon",
				Details:    err.Error(),
			}
			return
		}

		service, err := gopatreon.NewService(client)
		if err != nil {
			errCh <- &apierror.Error{
				StatusCode: http.StatusInternalServerError,
				Message:    errors.New("unable to create Patreon client").Error(),
				Details:    err.Error(),
			}
			return
		}

		userName, err := service.AuthenticateUser()
		if err != nil {
			errCh <- &apierror.Error{
				StatusCode: http.StatusForbidden,
				Message:    errors.New("unable to authenticate user").Error(),
				Details:    err.Error(),
			}
			return
		}

		baseURL := "http://localhost:3000/welcome"
		v := url.Values{}
		v.Set("name", userName)
		redirectURL := baseURL + "?" + v.Encode()

		newReq, err := http.NewRequestWithContext(r.Context(), "GET", redirectURL, nil)
		if err != nil {
			errCh <- &apierror.Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "unable to create redirect to welcome page",
				Details:    err.Error(),
			}
			return
		}

		http.Redirect(w, newReq, redirectURL, http.StatusTemporaryRedirect)
		errCh <- nil
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}
