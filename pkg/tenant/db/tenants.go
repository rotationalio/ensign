package db

import (
	"context"
	"time"
)

type Tenant struct {
	Created  time.Time
	Modified time.Time
}

func CreateTenant(ctx context.Context, tenant *Tenant) error {
	mu.RLock()
	defer mu.RUnlock()

	if !connected() {
		return ErrNotConnected
	}

	return nil
}
