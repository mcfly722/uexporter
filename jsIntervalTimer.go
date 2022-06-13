package main

import (
	"time"

	"github.com/mcfly722/goPackages/context"

	"github.com/dop251/goja"
)

type intervalTimer struct {
	jsScheduler   *jsScheduler
	jsRuntime     *JSRuntime
	handler       *goja.Callable
	intervalMS    int64
	startSpreadMS int64
	terminate     chan bool
}

func newIntervalTimer(jsScheduler *jsScheduler, jsRuntime *JSRuntime, handler *goja.Callable, intervalMS int64, startSpreadMS int64) timer {
	return &intervalTimer{
		jsScheduler:   jsScheduler,
		jsRuntime:     jsRuntime,
		handler:       handler,
		intervalMS:    intervalMS,
		startSpreadMS: startSpreadMS,
		terminate:     make(chan bool),
	}
}

// Go ...
func (timer *intervalTimer) Go(current context.Context) {

	duration := time.Duration(timer.startSpreadMS) * time.Millisecond

loop:
	for {
		select {
		case <-time.After(duration):
			timer.jsRuntime.CallHandler(timer.handler)
			duration = time.Duration(timer.intervalMS) * time.Millisecond

			break
		case <-timer.terminate:
			break loop
		case <-current.OnDone():
			break loop
		}
	}
}

// Dispose ...
func (timer *intervalTimer) Dispose() {
	timer.jsScheduler.deleteExistingTimer(timer)
}

func (timer *intervalTimer) Terminate() {
	timer.terminate <- true
}
