package debouncer

import (
	"errors"
	"sync"
	"time"
)

const (
	// ErrorTypeIncorrectSendSignal if you call SendSignal on something you configured WithAny
	ErrorTypeIncorrectSendSignal = "you are using sendsignal with something setup withany"
	// ErrorTypeIncorrectSendSignalWithAny  if you call SendSignalWithAny on something you configured WithTriggered
	ErrorTypeIncorrectSendSignalWithAny = "you are using sendsignalwithany with something setup withtriggered"
)

// Debouncer main struct for debouncer package
type Debouncer struct {
	timeDuration     time.Duration
	timer            *time.Timer
	triggeredFunc    func()
	triggeredAnyFunc func(any)
	isAny            bool
	mu               sync.Mutex
	done             chan struct{}
}

// New creates a new instance of debouncer. Each instance of debouncer works independent, concurrency with different wait duration.
func New(duration time.Duration) *Debouncer {
	return &Debouncer{timeDuration: duration, triggeredFunc: func() {}, triggeredAnyFunc: func(any) {}}
}

// WithTriggered attached a triggered function to debouncer instance and return the same instance of debouncer to use.
func (d *Debouncer) WithTriggered(triggeredFunc func()) *Debouncer {
	d.triggeredFunc = triggeredFunc
	d.isAny = false
	return d
}

// WithAny attached a triggered function to debouncer instance and return the same instance of debouncer to use.
func (d *Debouncer) WithAny(triggeredFunc func(any)) *Debouncer {
	d.triggeredAnyFunc = triggeredFunc
	d.isAny = true
	return d
}

// SendSignal makes an action that notifies to invoke the triggered function after a wait duration.
func (d *Debouncer) SendSignal() (err error) {
	if d.isAny {
		return errors.New(ErrorTypeIncorrectSendSignalWithAny)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.Cancel()
	d.timer = time.AfterFunc(d.timeDuration, func() {
		d.triggeredFunc()
		if d.done != nil {
			close(d.done)
		}
		d.done = make(chan struct{})
	})
	return nil
}

// SendSignalWithData makes an action that notifies to invoke the triggered function after a wait duration.
func (d *Debouncer) SendSignalWithData(anyVar any) (err error) {
	if !d.isAny {
		return errors.New(ErrorTypeIncorrectSendSignal)
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Cancel()
	d.timer = time.AfterFunc(d.timeDuration, func() {
		d.triggeredAnyFunc(anyVar)
		if d.done != nil {
			close(d.done)
		}
		d.done = make(chan struct{})
	})
	return nil
}

// Do run the signalFunc() and call SendSignal() after all. The signalFunc() and SendSignal() function run sequentially.
func (d *Debouncer) Do(signalFunc func()) {
	signalFunc()
	d.SendSignal()
}

// DoAny run the signalFunc(any) and call SendSignalWithData(any) after all. The signalFunc() and SendSignal() function run sequentially.
func (d *Debouncer) DoAny(signalFunc func(any), anyVar any) {
	signalFunc(anyVar)
	d.SendSignalWithData(anyVar)
}

// Cancel the timer from the last function SendSignal(). The scheduled triggered function is cancelled and doesn't invoke.
func (d *Debouncer) Cancel() {
	if d.timer != nil {
		d.timer.Stop()
	}
}

// UpdateTriggeredFunc replaces triggered function.
func (d *Debouncer) UpdateTriggeredFunc(newTriggeredFunc func()) {
	d.triggeredFunc = newTriggeredFunc
}

// UpdateAnyFunc replaces triggered function.
func (d *Debouncer) UpdateAnyFunc(newTriggeredFunc func(any)) {
	d.triggeredAnyFunc = newTriggeredFunc
}

// UpdateTimeDuration replaces the waiting time duration. You need to call a SendSignal() again to trigger a new timer with a new waiting time duration.
func (d *Debouncer) UpdateTimeDuration(newTimeDuration time.Duration) {
	d.timeDuration = newTimeDuration
}

// Done returns a receive-only channel to notify the caller when the triggered func has been executed.
func (d *Debouncer) Done() <-chan struct{} {
	if d.done == nil {
		d.done = make(chan struct{})
	}
	return d.done
}
