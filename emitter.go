package saga

import (
	"fmt"
	"sync"
)

type Emitter struct {
	events map[string]chan interface{}
	// https://go.dev/doc/effective_go#data
	// sync.Mutex does not have an explicit constructor or Init method. Instead, the zero value for a sync.Mutex is defined to be an unlocked mutex.
	mu sync.Mutex
}

func NewEmitter() *Emitter {
	return &Emitter{
		events: make(map[string]chan interface{}),
	}
}

func (e *Emitter) on(event string) chan interface{} {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.events[event]; !ok {
		e.events[event] = make(chan interface{})
	}

	return e.events[event]
}

func (e *Emitter) On(event Command, handler func(CommandHandler)) {
	CommandValidator{}.Validate(event)

	ch := e.on(event)
	go func() {
		for data := range ch {
			handler(data.(CommandHandler))
		}
		fmt.Printf("Channel %s closed. Exiting goroutine.\n", event)
	}()
}

func (e *Emitter) Off(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if event == "*" {
		for _, ch := range e.events {
			close(ch)
		}
		e.events = make(map[string]chan interface{})
	} else {
		if ch, ok := e.events[event]; ok {
			close(ch)
			delete(e.events, event)
		}
	}
}

func (e *Emitter) Emit(event string, data interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ch, ok := e.events[event]; ok {
		ch <- data
	}
}
