package api

import(
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
)

type (
	Welcome struct{
		OAuth2Config *oauth2.Config
	}
)

func (wel *Welcome) Handler(w http.ResponseWriter, r *http.Request) error {
	tok, err := wel.OAuth2Config.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Fatal(err)
	}

	patreonClient := patreon.NewClient(wel.OAuth2Config.Client(r.Context(), tok))
	user, err := patreonClient.FetchUser()
	if err != nil {
		return err
	}
	fmt.Println(user)

	// temporary static data
	welcomeMsg := map[string]string{
		"message": "welcome! you're logged in",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(welcomeMsg)
	if err != nil {
		return err
	}
	return nil
}
