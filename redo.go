package redo

import "time"
import "fmt"

// Func represents the function to pass to the Redo object
type Func func(stop func()) error

// ErrMaxRetryReached indicates that max retry has been reached
var ErrMaxRetryReached = fmt.Errorf("max retry reached")

// Redo defines a structure that provides the ability
// to run a function continuously as long as the function
// returns an error.
type Redo struct {
	maxRetries int
	retryDelay time.Duration
	stop       bool
	LastErr    error
}

// NewRedo creates a new Redo instance.
// Setting maxRetries to -1 will cause Redo to run the function forever.
func NewRedo(maxRetries int, retryDelay time.Duration) *Redo {
	return &Redo{
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Stop redoing
func (r *Redo) Stop() {
	r.stop = true
}

// Do runs a function. It will continuously retry the function
// if it returns errs abd will only stop if max retries is exceeded
// or explicitly stopped using the stop function passed to the
// running function or the object's Stop function. ErrMaxRetryReached is
// returned if the max retries has reached. Check the LastErr object field
// for the last error returned by the function.
func (r *Redo) Do(f Func) error {
	retryCount := 0
	for !r.stop {

		retryCount++

		if r.maxRetries > -1 && retryCount > r.maxRetries {
			return ErrMaxRetryReached
		}

		r.LastErr = f(r.Stop)
		if r.LastErr == nil {
			break
		}

		time.Sleep(r.retryDelay)
	}
	return r.LastErr
}
