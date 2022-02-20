package api

import (
	"net/http"

	"golang.org/x/oauth2"
)

type (
	Login struct {
		OAuth2Config *oauth2.Config
	}
)

func (l *Login) Handler(w http.ResponseWriter, r *http.Request) error {
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
		url := l.OAuth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		errCh <- nil
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}
