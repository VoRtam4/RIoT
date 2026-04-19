package events

import (
	"fmt"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type Subscription struct {
	ID      string
	Types   map[EventType]bool
	Filter  func(Event) bool
	Channel chan Event
}

type Bus struct {
	mutex sync.RWMutex
	subs  map[string]*Subscription
}

var (
	globalBus     *Bus
	globalBusOnce sync.Once
)

func GetEventBus() *Bus {
	globalBusOnce.Do(func() {
		globalBus = &Bus{
			subs: map[string]*Subscription{},
		}
	})
	return globalBus
}

func (b *Bus) Subscribe(eventTypes []EventType, filter func(Event) bool, bufferSize int) *Subscription {
	if bufferSize <= 0 {
		bufferSize = 16
	}
	typesMap := map[EventType]bool{}
	for _, eventType := range eventTypes {
		typesMap[eventType] = true
	}
	subscriptionID := fmt.Sprintf("sub-%s-%s", sharedUtils.GenerateRandomAlphanumericString(8), time.Now().UTC().Format(time.RFC3339Nano))
	subscription := &Subscription{
		ID:      subscriptionID,
		Types:   typesMap,
		Filter:  filter,
		Channel: make(chan Event, bufferSize),
	}
	b.mutex.Lock()
	b.subs[subscriptionID] = subscription
	b.mutex.Unlock()
	return subscription
}

func (b *Bus) Unsubscribe(subscriptionID string) {
	b.mutex.Lock()
	subscription, exists := b.subs[subscriptionID]
	if !exists {
		return
	}
	delete(b.subs, subscriptionID)
	b.mutex.Unlock()
	close(subscription.Channel)
}

func (b *Bus) Publish(eventType EventType, payload any) {
	event := Event{
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, subscription := range b.subs {
		if !subscription.Types[eventType] {
			continue
		}
		if subscription.Filter != nil && !subscription.Filter(event) {
			continue
		}
		select {
		case subscription.Channel <- event:
		default:
		}
	}
}
