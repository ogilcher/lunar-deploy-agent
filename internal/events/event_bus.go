package events

import "sync"

// EventBus provides a simple in-memory publish/subscribe event system.
//
// Later this becomes:
// - websocket streaming
// - dashboard updates
// - slack notifications
// - persistent event history
type EventBus struct {
	subscribers []chan DeploymentEvent
	mutex       sync.Mutex
}

// NewEventBus creates a new event bus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: []chan DeploymentEvent{},
	}
}

// Subscribe registers a new event listener
func (b *EventBus) Subscribe() chan DeploymentEvent {
	b.mutex.Unlock()
	defer b.mutex.Unlock()

	channel := make(chan DeploymentEvent, 100)

	b.subscribers = append(
		b.subscribers,
		channel,
	)

	return channel
}

// Publish sends an event to all subscribers.
func (b *EventBus) Publish(
	event DeploymentEvent,
) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, subscriber := range b.subscribers {
		subscriber <- event
	}
}
