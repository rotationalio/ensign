package api

import (
	"context"
)

// A Reauthenticator generates new access and refresh pair given a valid refresh token.
type Reauthenticator interface {
	Refresh(context.Context, *RefreshRequest) (*LoginReply, error)
}
