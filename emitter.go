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

// Custom constrainti
type PayloadEventMap interface {
	PaymentsNotifyClientPayload | SocialUnblockChatPayload
}
type EventHandler struct {
	Channel *EventsConsumeChannel `json:"channel"`
	Payload PayloadEvent          `json:"payload"`
}

// Custom constrainti
type PayloadType interface {
	SocialUnblockChatPayload | PaymentsNotifyClientPayload
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

//func (e *Emitter[T, U]) Off(event U) {
//	e.mu.Lock()
//	defer e.mu.Unlock()
//
//	if event == "*" {
//		for _, ch := range e.events {
//			close(ch)
//		}
//		e.events = make(map[string]chan T)
//	} else {
//		if ch, ok := e.events[event]; ok {
//			close(ch)
//			delete(e.events, event)
//		}
//	}
//}

func (e *Emitter[T, U]) Emit(event U, data T) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ch, ok := e.events[event]; ok {
		ch <- data
	}
}

func Pair[K comparable, V any](key K, value V) struct {
	Key   K
	Value V
} {
	return struct {
		Key   K
		Value V
	}{
		Key:   key,
		Value: value,
	}
}

type CommandMap[T any] map[string]func(T)

func RegisterCommand[T any](cm CommandMap[T], name string, handler func(T)) {
	cm[name] = handler
}

func ExecuteCommand[T any](cm CommandMap[T], name string, data T) {
	handler, ok := cm[name]
	if !ok {
		// Handle unknown command
		return
	}
	handler(data)
}

type Stringable interface {
	String() string
}

// Interface Constraints:
func ToString[T Stringable](x T) string {
	return x.String()
}

// Composite Constraints:
func Add[T int | float64](a, b T) T {
	return a + b
}

// Custom constrainti
type SignedInteger interface {
	int | int8 | int16 | int32 | int64
}

func Abs[T SignedInteger](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func foo() {
	a := Abs(-42)
	commandEmitter := newEmitter[CommandHandler, SagaStepCommand]()
	commandEmitter.On(NewUserSetRolesToRoomsCommand, func(handler CommandHandler) {
		payload := handler.Payload
		fmt.Println(payload)
	})

	eventEmitter := newEmitter[EventHandler, MicroserviceEvent]()
	eventEmitter.On(PaymentsNotifyClientEvent, func(handler EventHandler) {
		payload := handler.Payload.(PaymentsNotifyClientPayload)
		fmt.Println(payload)
	})
	eventEmitter.On(SocialUnblockChatEvent, func(handler EventHandler) {
		// Handle EventHandler specifically
		payload := handler.Payload
		fmt.Println(payload)
	})
	// examples with Pair
	pair := Pair[string, int]("key", 1)
	fmt.Println(pair)

	// Create a pair of an integer key and a string value
	intStringPair := Pair(42, "the answer")
	fmt.Println(intStringPair) // Output: {42 the answer}

	// Use the pair as a map entry
	myMap := map[int]string{intStringPair.Key: intStringPair.Value}
	fmt.Println(myMap) // Output: map[42:the answer]
	//
	cm := make(CommandMap[string])
	RegisterCommand(cm, "greet", func(name string) {
		fmt.Println("Hello,", name)
	})

	ExecuteCommand(cm, "greet", "Alice")
}
