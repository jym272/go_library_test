package saga

import (
	"fmt"
)

type EventsConsumeChannel struct {
	*ConsumeChannel
}

func (m *EventsConsumeChannel) AckMessage() {
	err := m.channel.Ack(m.msg.DeliveryTag, false)
	if err != nil {
		fmt.Println("error acknowledging message: %w", err)
	}
}
