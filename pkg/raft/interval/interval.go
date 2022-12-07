package interval

import (
	"math/rand"
	"time"
)

//===========================================================================
// Interval Interface
//===========================================================================

// Interval is an interface that specifies the behavior of time based signals.
// An interval dispatches a single signal to its channel, C at the delay specified by
// the schedule, which can be either fixed or stochastic. Fixed intervals resechedule
// themselves for a fixed delay after all the signal has been dispatched. Stochastic
// intervals select a random delay in a configured range to schedule the next event.
//
// Interval objects can be started and stopped. On start, the interval
// schedules the next event after the delay returned by GetDelay(). On stop
// no events will be dispatched by the handler. Intervals can be
// interrupted which resets the timer to a new delay. Timer state (running or
// not running) can be determined by the Running() method.
//
// Intervals are not thread safe and should only be used from within a single thread.
type Interval interface {
	Start() bool     // start the interval to periodically call its function
	Stop() bool      // stop the interval, the function will not be called
	Interrupt() bool // interrupt the interval, setting it to the next period
	Running() bool   // whether or not the interval is running
}

// NewFixed creates and initializes a new fixed interval.
func NewFixed(delay time.Duration) *FixedInterval {
	fixed := &FixedInterval{
		interval: interval{
			C: make(chan time.Time, 1),
		},
		delay: delay,
	}

	fixed.interval.timer = fixed
	return fixed
}

// NewRandom creates and initializes a new uniform random interval inside of a delay range.
func NewRandom(minDelay, maxDelay time.Duration) *RandomInterval {
	// Give the channel a 1-element time buffer.
	// If the client falls behind while reading, we drop ticks
	// on the floor until the client catches up.
	random := &RandomInterval{
		interval: interval{
			C: make(chan time.Time, 1),
		},
		minDelay: int64(minDelay),
		maxDelay: int64(maxDelay),
	}

	random.interval.timer = random
	return random
}

// NewJitter creates and initializes a new normal random interval with mean and stddev.
func NewJitter(meanDelay, stddevDelay time.Duration) *JitterInterval {
	// Give the channel a 1-element time buffer.
	// If the client falls behind while reading, we drop ticks
	// on the floor until the client catches up.
	random := &JitterInterval{
		interval: interval{
			C: make(chan time.Time, 1),
		},
		meanDelay:   float64(meanDelay),
		stddevDelay: float64(stddevDelay),
	}

	random.interval.timer = random
	return random
}

//===========================================================================
// FixedInterval Declaration
//===========================================================================

// FixedInterval sends a timestamp on its channel at a fixed interval specified by delay.
type FixedInterval struct {
	interval
	delay time.Duration
}

var _ Interval = &FixedInterval{}

func (t *FixedInterval) Initialized() bool {
	return t.delay > 0
}

// GetDelay returns the fixed interval duration.
func (t *FixedInterval) GetDelay() time.Duration {
	return t.delay
}

//===========================================================================
// RandomInterval Declaration
//===========================================================================

// RandomInterval sends a timestamp on its channel at a uniform random interval between
// the minimum and maximum delays. The maximum delay must be greater than the minimum.
type RandomInterval struct {
	interval
	minDelay int64
	maxDelay int64
}

var _ Interval = &RandomInterval{}

func (t *RandomInterval) Initialized() bool {
	return t.maxDelay-t.minDelay > 0
}

// GetDelay returns a random integer in the range (minDelay, maxDelay) on
// every request for the delay, causing jitter so that no timeout occurs at
// the same time.
func (t *RandomInterval) GetDelay() time.Duration {
	return time.Duration(rand.Int63n(t.maxDelay-t.minDelay) + t.minDelay)
}

//===========================================================================
// JitterInterval Declaration
//===========================================================================

// JitterInterval sends a timestamp on its channel at a normally distributed random
// interval with a mean and standard deviation delays. The mean and the stddev must be
// greater than 0. When sampling the distribution, the interval tries 7 times to get a
// non-zero delay, otherwise it defaults to the mean. It's important to choose a
// distribution that is unlikely to sample values less than or equal to zero.
type JitterInterval struct {
	interval
	meanDelay   float64
	stddevDelay float64
}

func (t *JitterInterval) Initialized() bool {
	return t.meanDelay > 0 && t.stddevDelay > 0
}

// GetDelay returns a random delay with a normal distribution of mean delay and a
// standard deviation of stddev delay. This method tries 7 times to return a non-zero,
// non-negative delay then defaults to returning the mean.
func (t *JitterInterval) GetDelay() time.Duration {
	for i := 0; i < 7; i++ {
		if samp := time.Duration(rand.NormFloat64()*t.stddevDelay + t.meanDelay); samp > 0 {
			return samp
		}
	}
	return time.Duration(t.meanDelay)
}

//===========================================================================
// Base Interval
//===========================================================================

// All structs that embed interval must implement the Timer interface.
type Timer interface {
	Initialized() bool
	GetDelay() time.Duration
}

// An embedded interval that implements the Interval interface for most Intervals so
// long as the struct embedding interval implements the Timer interface.
type interval struct {
	C       chan time.Time
	timer   Timer
	running bool
	stop    chan struct{}
}

// Start the interval to periodically issue events. Returns true if the
// ticker gets started, false if it's already started or uninitialized.
func (t *interval) Start() bool {
	// If the delay is 0 or negative there is no reason to start the ticker
	// Should not be able to start an already running interval
	if t.timer == nil || !t.timer.Initialized() || t.running {
		return false
	}

	stop := make(chan struct{}, 1)
	t.stop = stop

	go t.loop(stop)
	t.running = true
	return true
}

// Stop the interval so that no more events are dispatched. Returns true if
// the call stops the interval, false if already expired or never started.
func (t *interval) Stop() bool {
	if !t.running {
		return false
	}

	close(t.stop)
	t.stop = nil
	t.running = false
	return true
}

// Interrupt the current interval, stopping and starting it again. Returns
// true if the interval was running and is successfully reset, false if the
// ticker was stopped or uninitialized.
func (t *interval) Interrupt() bool {
	if !t.running {
		return false
	}

	// Stop the current loop
	close(t.stop)

	// Create another loop
	stop := make(chan struct{}, 1)
	t.stop = stop

	go t.loop(stop)
	t.running = true
	return true
}

// Running returns true if the timer exists and false otherwise.
func (t *interval) Running() bool {
	return t.running
}

// A go routine that sends timestamps on the internal channel.
func (t *interval) loop(stop <-chan struct{}) {
	for {
		wait := time.After(t.timer.GetDelay())
		select {
		case ts := <-wait:
			select {
			case t.C <- ts:
				continue
			default:
			}
		case <-stop:
			return
		}
	}
}
