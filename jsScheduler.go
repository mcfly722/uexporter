package main

import (
	"fmt"
	"sync"

	"github.com/dop251/goja"
	"github.com/mcfly722/goPackages/context"
	"github.com/mcfly722/goPackages/logger"
)

type timer interface {
	context.ContextedInstance
	Terminate()
}

type jsScheduler struct {
	name          string
	logger        *logger.Logger
	activeTimers  map[int64]timer
	timersCounter int64

	ready sync.Mutex
}

// NewJSScheduler ...
func NewJSScheduler(name string, logger *logger.Logger) API {
	return &jsScheduler{
		name:          fmt.Sprintf("scheduler %v", name),
		logger:        logger,
		activeTimers:  make(map[int64]timer),
		timersCounter: 0,
	}
}

func (jsScheduler *jsScheduler) addNewTimer(timer timer) int64 {
	jsScheduler.ready.Lock()
	defer jsScheduler.ready.Unlock()

	timerID := jsScheduler.timersCounter
	jsScheduler.activeTimers[timerID] = timer

	jsScheduler.logger.LogEvent(logger.EventTypeTrace, jsScheduler.name, fmt.Sprintf("%v timer created", timerID))

	jsScheduler.timersCounter++
	return timerID
}

func (jsScheduler *jsScheduler) deleteExistingTimer(timer timer) {
	jsScheduler.ready.Lock()
	defer jsScheduler.ready.Unlock()
	timerID := int64(-1)

	for id, timerPointer := range jsScheduler.activeTimers {
		if timer == timerPointer {
			timerID = id
		}
	}

	if timerID != -1 {
		delete(jsScheduler.activeTimers, timerID)
		jsScheduler.logger.LogEvent(logger.EventTypeTrace, jsScheduler.name, fmt.Sprintf("%v timer terminated", timerID))
	}
}

func (jsScheduler *jsScheduler) tryToDeleteTimer(timerID int64) error {
	jsScheduler.ready.Lock()
	defer jsScheduler.ready.Unlock()

	if timer, ok := jsScheduler.activeTimers[timerID]; ok {
		timer.Terminate()
		return nil
	}
	return fmt.Errorf("timer with id=%v does not exist", timerID)
}

func (jsScheduler *jsScheduler) Init(current context.Context, jsRuntime *JSRuntime) {

	setInterval := func(handler goja.Callable, intervalMS int64, startSpreadMS int64) int64 {
		timer := newIntervalTimer(jsScheduler, jsRuntime, &handler, intervalMS, startSpreadMS)
		timerID := jsScheduler.addNewTimer(timer)
		current.NewContextFor(timer)
		return timerID
	}

	clearInterval := func(timerId int64) {
		if err := jsScheduler.tryToDeleteTimer(timerId); err != nil {
			jsRuntime.Throw(fmt.Sprintf("clearInterval(%v): %v", timerId, err))
		}

	}

	jsRuntime.VM.Set("setInterval", setInterval)
	jsRuntime.VM.Set("clearInterval", clearInterval)

}
