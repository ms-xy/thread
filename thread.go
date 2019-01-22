/*
Package thread provides a "Thread"-like convenience wrapper around goroutines.
*/
package thread

import (
	"errors"
	"sync"
)

// State type determines a Thread's execution status
type State uint8

const (
	RUNNING State = iota
	STOPPING
	STOPPED
)

var (
	ErrAlreadyInitialized = errors.New("Thread has already been initialized")
	ErrAlreadyStarted     = errors.New("Thread has already been started")
	ErrMalfunction        = errors.New("Thread state is broken")
)

// The Thread struct is neither a kernel nor a user thread implementation.
// All it actually does is executing a goroutine and providing means to start
// and stop it. Call it thread-like if you like.
type Thread struct {
	mutex        sync.Mutex
	initialized  bool
	state        State
	stopRunnable chan bool
	waitThread   chan bool
	runnable     Runnable
}

// Runnable is a simple interface describing a minimalistic runnable type
// consisting of a main loop and using a for-select to determine when a
// stop is expected:
//
//   for {
//     select {
//     case <-stop:
//       return nil
//     default:
//       // do work
//     }
//   }
//
// See the package example for details how to implement such a runnable.
type Runnable interface {
	Run(stop chan bool) error
}

// New creates a new Thread and initializes it with the given Runnable.
// Must be started separately using Thread.Start()
func New(runnable Runnable) *Thread {
	return (&Thread{}).Init(runnable)
}

// Init initializes the Thread with the given Runnable.
// Panics with ErrAlreadyInitialized if it has been initialized before.
func (t *Thread) Init(runnable Runnable) *Thread {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	// check state, if already initialized, panic
	if t.initialized {
		panic(ErrAlreadyInitialized)
	}
	// set initial field values
	t.initialized = true
	t.state = STOPPED
	t.runnable = runnable
	return t
}

// Start starts the Thread in a new goroutine and initializes its signal channels.
func (t *Thread) Start() {
	// check if already running
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.state != STOPPED {
		return
	}
	// setup signal channels and update state to running
	t.stopRunnable = make(chan bool)
	t.waitThread = make(chan bool)
	t.state = RUNNING
	// launch new goroutine
	go t.run()
}

// Internal helper function for running then cleaning up
func (t *Thread) run() {
	defer func() {
		t.mutex.Lock()
		defer t.mutex.Unlock()
		// in case we haven't been stopped, the channel is still open, so close it
		if t.state != STOPPING {
			close(t.stopRunnable)
		}
		// indicate state change and close wait thread in case anyone is listening
		t.state = STOPPED
		close(t.waitThread)
	}()
	// run child
	t.runnable.Run(t.stopRunnable)
}

// Stop the Thread by signaling the Runnable to stop, effectively resulting in the target goroutine to exit.
// To wait for the Thread to finish use Thread.Join().
func (t *Thread) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	// check state, stopping twice is useless, so simply return
	if t.state != RUNNING {
		return
	}
	t.state = STOPPING
	// signal the runnable to stop
	close(t.stopRunnable)
}

// Join blocks until the Thread terminates.
func (t *Thread) Join() {
	// wait until runnable has exited
	<-t.waitThread
}
