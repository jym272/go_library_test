package saga

import (
	"fmt"
	"sync"
)

type Emitter[T any, U comparable] struct {
	events map[U]chan T
	// https://go.dev/doc/effective_go#data
	// sync.Mutex does not have an explicit constructor or Init method. Instead, the zero value for a sync.Mutex is defined to be an unlocked mutex.
	mu sync.Mutex
}

func newEmitter[T any, U comparable]() *Emitter[T, U] {
	return &Emitter[T, U]{
		events: make(map[U]chan T),
	}
}

func (e *Emitter[T, U]) on(event U) chan T {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.events[event]; !ok {
		e.events[event] = make(chan T)
	}

	return e.events[event]
}

type EventHandler struct {
	Channel *EventsConsumeChannel `json:"channel"`
	Payload PayloadEvent          `json:"payload"`
}

func (e *Emitter[T, U]) On(event U, handler func(T)) {
	//CommandValidator{}.Validate(event)

	ch := e.on(event)
	go func() {
		for data := range ch {
			handler(data)
		}
		fmt.Printf("Channel %s closed. Exiting goroutine.\n", event)
	}()
}

func (e *Emitter[T, U]) Emit(event U, data T) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ch, ok := e.events[event]; ok {
		ch <- data
	}
}
