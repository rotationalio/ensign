package contexts

import "errors"

var (
	ErrNoClaimsInContext = errors.New("no claims available in context")
	ErrNotAuthorized     = errors.New("claims do not have required permission")
)
