package events

import (
	"testing"
	"time"
)

func TestEventBusPublishDeliversEvent(t *testing.T) {
	bus := NewEventBus()
	subscriber := bus.Subscribe()
	defer bus.Unsubscribe(subscriber)

	expected := DeploymentEvent{
		Type:       "deployment_started",
		Deployment: "portfolio",
		Message:    "deployment started",
	}

	bus.Publish(expected)

	select {
	case actual := <-subscriber:
		if actual.Type != expected.Type {
			t.Fatalf(
				"expected type %q, got %q",
				expected.Type,
				actual.Type,
			)
		}

		if actual.Deployment != expected.Deployment {
			t.Fatalf(
				"expected deployment %q, got %q",
				expected.Deployment,
				actual.Deployment,
			)
		}

	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for event")
	}
}

func TestEventBusUnsubscribeStopsFutureEvents(t *testing.T) {
	bus := NewEventBus()
	subscriber := bus.Subscribe()

	bus.Unsubscribe(subscriber)

	bus.Publish(DeploymentEvent{
		Type: "deployment_started",
	})

	select {
	case event := <-subscriber:
		t.Fatalf(
			"unexpected event after unsubscribe: %+v",
			event,
		)

	case <-time.After(25 * time.Millisecond):
	}
}

func TestEventBusSlowSubscriberDoesNotBlockPublish(t *testing.T) {
	bus := NewEventBus()
	subscriber := bus.Subscribe()
	defer bus.Unsubscribe(subscriber)

	for i := 0; i < subscriberBufferSize; i++ {
		bus.Publish(DeploymentEvent{
			Type: "buffer_fill",
		})
	}

	done := make(chan struct{})

	go func() {
		bus.Publish(DeploymentEvent{
			Type: "should_not_block",
		})

		close(done)
	}()

	select {
	case <-done:

	case <-time.After(100 * time.Millisecond):
		t.Fatal("publish blocked on a slow subscriber")
	}
}
