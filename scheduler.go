package main

import (
	"fmt"
	"sync"
)

// Scheduler ...
type Scheduler struct {
	counter     uint64
	pluginsPath string
	ready       sync.Mutex
}

// NewScheduler ...
func NewScheduler(pluginsPath string) *Scheduler {
	scheduler := &Scheduler{
		counter:     0,
		pluginsPath: pluginsPath,
	}

	return scheduler
}

// ToHTML ...
func (scheduler *Scheduler) ToHTML() string {
	scheduler.ready.Lock()

	out := fmt.Sprintf("counter:%v<br>", scheduler.counter)
	out += fmt.Sprintf("pluginsPath:%v<br>", scheduler.pluginsPath)

	scheduler.counter++

	scheduler.ready.Unlock()
	return out
}
