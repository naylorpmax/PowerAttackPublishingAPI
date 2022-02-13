package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"

	"github.com/naylorpmax/homebrew-users-api/pkg/middleware"
)

type (
	Welcome struct {
		OAuth2Config *oauth2.Config
	}

	Patreon struct {
		client *patreon.Client
	}
)

const (
	CreatorUserID      = "12794096"
	CampaignID         = "1976402"
	MinUserAmountCents = 500
)

func NewPatreon(r *http.Request, oauth2Config *oauth2.Config) (*Patreon, error) {
	code := r.FormValue("code")
	if code == "" {
		return nil, errors.New("redirect request does not contain OAuth2 code")
	}

	tok, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to create Patreon client: %v", err.Error())
	}

	client := oauth2Config.Client(r.Context(), tok)
	return &Patreon{client: patreon.NewClient(client)}, nil
}

func (wel *Welcome) Handler(w http.ResponseWriter, r *http.Request) error {
	patreonClient, err := NewPatreon(r, wel.OAuth2Config)
	if err != nil {
		return &middleware.Error{
			StatusCode: http.StatusForbidden,
			Message:    errors.New("unable to authenticate to Patreon").Error(),
			Details:    err.Error(),
		}
	}

	userName, err := patreonClient.authenticateUser()
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

func (p *Patreon) authenticateUser() (string, error) {
	user, err := p.client.FetchUser()
	if err != nil {
		return "", err
	}

	// hello creator!
	if p.getUserID(user) == CreatorUserID {
		return p.getUserFirstLastName(user), nil
	}

	// hello patron!
	pledges, err := p.client.FetchPledges(CampaignID)
	if err != nil {
		return "", err
	}
	if pledgeAmount := p.getPledgeAmount(pledges); pledgeAmount < MinUserAmountCents {
		return "", errors.New("patron level not high enough to access content")
	}
	if err = p.goodStanding(user, pledges); err != nil {
		return "", fmt.Errorf("user is not in good standing with this campaign: %v", err.Error())
	}
	return p.getUserFirstLastName(user), nil
}

func (p *Patreon) getUserID(user *patreon.UserResponse) string {
	return user.Data.ID
}

func (p *Patreon) getUserFirstLastName(user *patreon.UserResponse) string {
	return user.Data.Attributes.FirstName + " " + user.Data.Attributes.LastName
}

func (p *Patreon) getPledgeAmount(pledges *patreon.PledgeResponse) int {
	totalAmount := 0
	for _, pledge := range pledges.Data {
		totalAmount += pledge.Attributes.AmountCents
	}
	return totalAmount
}

func (p *Patreon) goodStanding(user *patreon.UserResponse, pledges *patreon.PledgeResponse) error {
	if user.Data.Attributes.IsSuspended {
		return errors.New("user is suspended")
	}
	if user.Data.Attributes.IsDeleted {
		return errors.New("user is deleted")
	}
	for _, pledge := range pledges.Data {
		if !pledge.Attributes.PatronPaysFees {
			return errors.New("user has unpaid fees")
		}
		if *pledge.Attributes.IsPaused {
			return errors.New("user is paused")
		}
	}
	return nil
}
