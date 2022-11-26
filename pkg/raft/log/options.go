package log

type Option func(l *Log) error

func WithStateMachine(sm StateMachine) Option {
	return func(l *Log) error {
		l.sm = sm
		return nil
	}
}

func WithSync(sync Sync) Option {
	return func(l *Log) error {
		l.sync = sync
		return nil
	}
}
