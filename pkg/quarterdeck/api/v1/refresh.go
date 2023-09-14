package api

import (
	"context"

	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
)

// Refresher implements the token Refresher interface by calling Quarterdeck.
type Refresher struct {
	client QuarterdeckClient
}

var _ tokens.Refresher = &Refresher{}

// NewRefresher returns a new Refresher that uses the Quarterdeck API to refresh tokens.
func NewRefresher(client QuarterdeckClient) *Refresher {
	return &Refresher{
		client: client,
	}
}

// Refresh an access token by calling Quarterdeck.
func (q *Refresher) Refresh(ctx context.Context, refresh string) (accessToken, refreshToken string, err error) {
	var rep *LoginReply
	if rep, err = q.client.Refresh(ctx, &RefreshRequest{
		RefreshToken: refreshToken,
	}); err != nil {
		return "", "", err
	}

	return rep.AccessToken, rep.RefreshToken, nil
}
