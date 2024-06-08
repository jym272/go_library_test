package saga

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CommandHandler struct {
	Channel *MicroserviceConsumeChannel `json:"channel"`
	Payload map[string]interface{}      `json:"payload"`
	SagaID  int                         `json:"sagaId"`
}

func microserviceConsumeCallback(msg *amqp.Delivery, channel *amqp.Channel, e *Emitter[CommandHandler, SagaStepCommand], queueName string) {
	if msg == nil {
		fmt.Println("NO MSG AVAILABLE")
		return
	}

	var currentStep SagaStep
	err := json.Unmarshal(msg.Body, &currentStep)
	if err != nil {
		fmt.Println("ERROR PARSING MSG", err)
		err := channel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}

	responseChannel := &MicroserviceConsumeChannel{
		step: currentStep,
		ConsumeChannel: &ConsumeChannel{
			channel:   channel,
			msg:       msg,
			queueName: queueName,
		},
	}

	e.Emit(currentStep.Command, CommandHandler{
		Channel: responseChannel,
		Payload: currentStep.PreviousPayload,
		SagaID:  currentStep.SagaID,
	})
}

// EventPayload is a generic container for microservice events
type EventPayload struct {
	EventType MicroserviceEvent
	Payload   interface{} // This will hold the event-specific data
}

// Helper function to create an EventPayload
func NewEventPayload(eventType MicroserviceEvent, payload interface{}) *EventPayload {
	return &EventPayload{
		EventType: eventType,
		Payload:   payload,
	}
}

// Example usage
func main() {
	// Create a TestImageEvent payload
	testImagePayload := TestImagePayload{Image: "https://example.com/image.jpg"}
	event := NewEventPayload(TestImageEvent, testImagePayload)

	// Accessing the payload data (type assertion)
	if event.EventType == TestImageEvent {
		imageData, ok := event.Payload.(TestImagePayload)
		if ok {
			// ... use imageData
		}
	}

	// ... similarly for other event types
}

type MicroserviceConsumeEvents struct {
	Payload interface{}
	Channel *EventsConsumeChannel
}

// eventCallback handles the consumption and processing of microservice events.
func eventCallback(msg *amqp.Delivery, channel *amqp.Channel, emitter *Emitter[EventHandler, MicroserviceEvent], queueName string) { // Assuming you have an Emitter type for event broadcasting
	if msg == nil {
		fmt.Println("Message not available")
		return
	}

	// Message parsing (with error handling and type assertion)
	var eventPayload EventPayload
	if err := json.Unmarshal(msg.Body, &eventPayload); err != nil {
		fmt.Printf("Error parsing message: %s\n", err)
		channel.Nack(msg.DeliveryTag, false, false) // Nack without requeue
		return
	}

	// Validate the event type for type safety
	if !isValidMicroserviceEvent(eventPayload.EventType) {
		fmt.Printf("Invalid event type: %s\n", eventPayload.EventType)
		channel.Nack(msg.DeliveryTag, false, false) // Nack without requeue
		return
	}

	// Emit the event to the appropriate channel
	eventTypeName := string(eventPayload.EventType)
	e.Emit(eventTypeName, eventPayload)

	// Acknowledge successful processing
	channel.Ack(msg.DeliveryTag, false)

	// Message parsing
	var payload interface{} // Use a generic interface to hold the payload initially
	err := json.Unmarshal(msg.Body, &payload)
	if err != nil {
		fmt.Println("Error parsing message:", err)
		channel.Nack(msg.DeliveryTag, false, false) // NACK without requeue
		return
	}

	// Extract the event key from headers
	eventKey, ok := findEventKey(msg.Headers)
	if !ok {
		fmt.Println("Invalid header value: no valid event key found")
		channel.Nack(msg.DeliveryTag, false, false)
		return
	}

	// Create the specific payload struct based on the event type
	var specificPayload interface{}
	switch eventKey {
	case TestImageEvent:
		specificPayload = &TestImagePayload{}
	case TestMintEvent:
		specificPayload = &TestMintPayload{}
	// ... add cases for other event types
	default:
		fmt.Println("Unsupported event type:", eventKey)
		channel.Nack(msg.DeliveryTag, false, false)
		return
	}

	// Unmarshal into the specific payload struct
	err = json.Unmarshal(msg.Body, specificPayload)
	if err != nil {
		fmt.Println("Error parsing specific payload:", err)
		channel.Nack(msg.DeliveryTag, false, false)
		return
	}

	responseChannel := &EventsConsumeChannel{
		channel:   channel,
		msg:       *msg,
		queueName: queueName,
	}

	// Emit the event with the typed payload and response channel
	emitter.Emit(eventKey, MicroserviceConsumeEvents{
		Payload: specificPayload,
		Channel: responseChannel,
	})
}
func eventCallback1(msg *amqp.Delivery, channel *amqp.Channel, e *Emitter, queueName string) {
	if msg == nil {
		fmt.Println("NO MSG AVAILABLE")
		return
	}

	var currentStep SagaStep
	err := json.Unmarshal(msg.Body, &currentStep)
	if err != nil {
		fmt.Println("ERROR PARSING MSG", err)
		err := channel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}

	responseChannel := &EventsConsumeChannel{
		ConsumeChannel: &ConsumeChannel{
			channel:   channel,
			msg:       msg,
			queueName: queueName,
		},
	}

	e.Emit(currentStep.Command, CommandHandler{
		Channel: responseChannel,
		Payload: currentStep.PreviousPayload,
		SagaID:  currentStep.SagaID,
	})
}

func consume[T any, U comparable](e *Emitter[T, U], queueName string, cb func(*amqp.Delivery, *amqp.Channel, *Emitter[T, U], string)) error {
	channel, err := getConsumeChannel()
	if err != nil {
		return err
	}

	channelQ, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range channelQ {
		cb(&msg, channel, e, queueName)
	}

	return nil
}

func PrintSlice[T any](s []T) {
	for _, v := range s {
		fmt.Println(v)
	}
}

func foo() {
	intSlice := []int{1, 2, 3}
	stringSlice := []string{"hello", "world"}

	PrintSlice[int](intSlice)
	PrintSlice[string](stringSlice)
}

func Min[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		panic("empty slice")
	}

	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min
}

func main1() {
	intSlice := []int{3, 1, 4, 2}
	floatSlice := []float64{3.14, 1.61, 4.67, 2.71}

	fmt.Println(Min[int](intSlice))       // Output: 1
	fmt.Println(Min[float64](floatSlice)) // Output: 1.61
}

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	lastIndex := len(s.data) - 1
	value := s.data[lastIndex]
	s.data = s.data[:lastIndex]
	return value, true
}

func (s *Stack[T]) Size() int {
	return len(s.data)
}

// https://medium.com/hprog99/mastering-generics-in-go-a-comprehensive-guide-4d05ec4b12b
func main4() {
	intStack := Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)

	fmt.Println(intStack.Pop())  // Output: 3, true
	fmt.Println(intStack.Size()) // Output: 2

	stringStack := Stack[string]{}
	stringStack.Push("hello")
	stringStack.Push("world")

	fmt.Println(stringStack.Pop())  // Output: world, true
	fmt.Println(stringStack.Size()) // Output: 1
}

func Map[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

type Person struct {
	Name string
	Age  int
}

func foo() {

	numbers := []int{1, 2, 3, 4, 5}
	squared := Map(numbers, func(x int) int { return x * x })
	fmt.Println(squared) // Output: [1 4 9 16 25]

	words := []string{"hello", "world", "Go"}
	uppercased := Map(words, strings.ToUpper)
	fmt.Println(uppercased) // Output: [HELLO WORLD GO]

	people := []Person{{"Alice", 30}, {"Bob", 25}}
	names := Map(people, func(p Person) string { return p.Name })
	fmt.Println(names) // Output: [Alice Bob]

	sendPayload(TextPayload{"Hello, world!"})
	sendPayload(ImagePayload{"https://example.com/image.jpg"})
}
