package patreon

import (
	"errors"
	"fmt"

	pat "github.com/naylorpmax/homebrew-users-api/pkg/client/patreon"
)

type (
	Patreon struct {
		Client pat.Client
	}
)

const (
	CreatorUserID      = "12794096"
	CampaignID         = "1976402"
	MinUserAmountCents = 500
)

func New(client pat.Client) (*Patreon, error) {
	return &Patreon{Client: client}, nil
}

func (p *Patreon) AuthenticateUser() (string, error) {
	user, err := p.Client.FetchUser()
	if err != nil {
		return "", fmt.Errorf("unable to fetch user: %v", err)
	}

	// hello creator!
	if user.ID == CreatorUserID {
		return user.FirstName + " " + user.LastName, nil
	}

	// hello patron!
	pledges, err := p.Client.FetchPledges(CampaignID)
	if err != nil {
		return "", fmt.Errorf("unable to fetch user's pledges: %v", err)
	}
	if pledgeAmount := getPledgeAmount(pledges); pledgeAmount < MinUserAmountCents {
		return "", errors.New("patron level not high enough to access content")
	}
	if err = goodStanding(user, pledges); err != nil {
		return "", fmt.Errorf("user is not in good standing with this campaign: %v", err.Error())
	}
	return user.FirstName + " " + user.LastName, nil
}

func getPledgeAmount(pledges []*pat.Pledge) int {
	totalAmount := 0
	for _, pledge := range pledges {
		totalAmount += pledge.AmountCents
	}
	return totalAmount
}

func goodStanding(user *pat.User, pledges []*pat.Pledge) error {
	if user.IsSuspended {
		return errors.New("user is suspended")
	}
	if user.IsDeleted {
		return errors.New("user is deleted")
	}
	if user.IsNuked {
		return errors.New("user is nuked")
	}
	for _, pledge := range pledges {
		if !pledge.PatronPaysFees {
			return errors.New("user has unpaid fees")
		}
		if pledge.IsPaused != nil && *pledge.IsPaused {
			return errors.New("user is paused")
		}
	}
	return nil
}
