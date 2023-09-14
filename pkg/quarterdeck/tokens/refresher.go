package tokens

import (
	"context"
)

// A Refresher generates new access and refresh pair given a valid refresh token.
type Refresher interface {
	Refresh(ctx context.Context, refresh string) (accessToken, refreshToken string, err error)
}
