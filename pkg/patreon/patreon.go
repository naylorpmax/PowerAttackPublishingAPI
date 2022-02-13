package patreon

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	mxpv "gopkg.in/mxpv/patreon-go.v1"
)

type (
	Patreon struct {
		client *mxpv.Client
	}
)

const (
	CreatorUserID      = "12794096"
	CampaignID         = "1976402"
	MinUserAmountCents = 500
)

// TODO: create interface + tests + consider making type to wrap mxpv user + pledge responses

func New(r *http.Request, oauth2Config *oauth2.Config) (*Patreon, error) {
	code := r.FormValue("code")
	if code == "" {
		return nil, errors.New("redirect request does not contain OAuth2 code")
	}

	tok, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to create Patreon client: %v", err.Error())
	}

	client := oauth2Config.Client(r.Context(), tok)
	return &Patreon{client: mxpv.NewClient(client)}, nil
}

func (p *Patreon) AuthenticateUser() (string, error) {
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

func (p *Patreon) getUserID(user *mxpv.UserResponse) string {
	return user.Data.ID
}

func (p *Patreon) getUserFirstLastName(user *mxpv.UserResponse) string {
	return user.Data.Attributes.FirstName + " " + user.Data.Attributes.LastName
}

func (p *Patreon) getPledgeAmount(pledges *mxpv.PledgeResponse) int {
	totalAmount := 0
	for _, pledge := range pledges.Data {
		totalAmount += pledge.Attributes.AmountCents
	}
	return totalAmount
}

func (p *Patreon) goodStanding(user *mxpv.UserResponse, pledges *mxpv.PledgeResponse) error {
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
