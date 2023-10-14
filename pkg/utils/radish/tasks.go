package radish

import "context"

type Task interface {
	Do(context.Context) error
}

type Func func(context.Context) error

func (f Func) Do(ctx context.Context) error {
	return f(ctx)
}
