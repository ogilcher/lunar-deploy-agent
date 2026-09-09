package events

import "sync"

const subscriberBufferSize = 100

// EventBus provides a simple in-memory publish/subscribe event system.
//
// Slow subscribers are isolated from publishers so deployment execution
// cannot be blocked by a disconnected or stalled event consumer.
type EventBus struct {
	subscribers map[chan DeploymentEvent]struct{}
	mutex       sync.RWMutex
}

// NewEventBus creates a new event bus instance.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[chan DeploymentEvent]struct{}),
	}
}

// Subscribe registers a new event listener.
func (b *EventBus) Subscribe() chan DeploymentEvent {
	channel := make(chan DeploymentEvent, subscriberBufferSize)

	b.mutex.Lock()
	b.subscribers[channel] = struct{}{}
	b.mutex.Unlock()

	return channel
}

// Unsubscribe removes an event listener from the bus.
//
// The channel is intentionally not closed. A publisher may already have a
// snapshot containing the channel, and closing it here could race with a send.
func (b *EventBus) Unsubscribe(channel chan DeploymentEvent) {
	b.mutex.Lock()
	delete(b.subscribers, channel)
	b.mutex.Unlock()
}

// Publish sends an event to all current subscribers.
//
// Delivery is best-effort. If a subscriber's buffer is full, the event is
// dropped for that subscriber rather than blocking deployment execution.
func (b *EventBus) Publish(event DeploymentEvent) {
	b.mutex.RLock()

	subscribers := make(
		[]chan DeploymentEvent,
		0,
		len(b.subscribers),
	)

	for subscriber := range b.subscribers {
		subscribers = append(subscribers, subscriber)
	}

	b.mutex.RUnlock()

	for _, subscriber := range subscribers {
		select {
		case subscriber <- event:
		default:
		}
	}
}
