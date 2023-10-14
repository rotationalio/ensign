package radish

import "context"

// Workers in the task manager handle Tasks which can hold state and other information
// needed by the task. You can also specify a simple function to execute by using the
// TaskFunc to create a Task to provide to the task manager.
type Task interface {
	Do(context.Context) error
}

type Func func(context.Context) error

func (f Func) Do(ctx context.Context) error {
	return f(ctx)
}
