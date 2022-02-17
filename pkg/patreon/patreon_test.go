package patreon_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pat "github.com/naylorpmax/homebrew-users-api/pkg/client/patreon"
	"github.com/naylorpmax/homebrew-users-api/pkg/patreon"
)

func TestPatreon_AuthenticateUser(t *testing.T) {
	type (
		in struct {
			fetchUser    func() (*pat.User, error)
			fetchPledges func(string) ([]*pat.Pledge, error)
		}
		exp struct {
			fullName string
			err      error
		}
	)
	cases := []struct {
		name     string
		input    in
		expected exp
	}{
		{
			name: "happy-patron",
			input: in{
				fetchUser: func() (*pat.User, error) {
					return &pat.User{
						ID: "123456",
						Attributes: pat.Attributes{
							FirstName:   "max",
							LastName:    "naylor",
							IsDeleted:   false,
							IsNuked:     false,
							IsSuspended: false,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*pat.Pledge, error) {
					isPaused := false
					return []*pat.Pledge{
						{
							AmountCents:    500,
							IsPaused:       &isPaused,
							PatronPaysFees: true,
						},
					}, nil
				},
			},
			expected: exp{
				fullName: "max naylor",
				err:      nil,
			},
		},
	}

	for _, test := range cases {
		client := mockClient{
			fetchUserFn:    test.input.fetchUser,
			fetchPledgesFn: test.input.fetchPledges,
		}
		patreonClient := patreon.Patreon{
			Client: &client,
		}

		actualFullName, actualErr := patreonClient.AuthenticateUser()

		assert.Equal(t, test.expected.fullName, actualFullName)
		assert.Equal(t, test.expected.err, actualErr)
	}
}

type mockClient struct {
	fetchUserFn    func() (*pat.User, error)
	fetchPledgesFn func(string) ([]*pat.Pledge, error)
}

func (m *mockClient) FetchUser() (*pat.User, error) {
	return m.fetchUserFn()
}

func (m *mockClient) FetchPledges(campaignID string) ([]*pat.Pledge, error) {
	return m.fetchPledgesFn(campaignID)
}
