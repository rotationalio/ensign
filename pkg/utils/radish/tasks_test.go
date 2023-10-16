package radish_test

import (
	"context"
	"fmt"
	"sync"
)

type TestTask struct {
	failUntil int
	attempts  int
	success   bool
	wg        *sync.WaitGroup
}

func (t *TestTask) Do(ctx context.Context) error {
	t.attempts++
	if t.attempts < t.failUntil {
		t.success = false
		return fmt.Errorf("task errored on attempt %d", t.attempts)
	}

	t.success = true
	t.wg.Done()
	return nil
}

func (t *TestTask) String() string {
	return "test task"
}
