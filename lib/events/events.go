package events

import (
	"fmt"
	"sort"
	"time"
)

// Events ...
type Events []Event

// Filter ... Filter EC2Events containing substr in Description
func (evs Events) Filter(states ...EventState) Events {
	events := Events{}

EVENTS:
	for _, ev := range evs {
		for _, state := range states {
			if ev.State == state {
				continue EVENTS // Skip and Next Event
			}
		}
		events = append(events, ev)
	}
	return events
}

// GetCloseEvent ... Get oldest EC2Event
func (evs Events) GetCloseEvent() Event {
	// Copy to temporary variable
	// Prevent to mutate self
	events := make(Events, cap(evs))
	copy(events, evs)

	// Sort as NotBefore date
	sort.Stable(events)

	return events[0]
}

// BeforeAll ...
func (evs Events) BeforeAll(d time.Time) bool {
	for _, ev := range evs {
		if ev.NotBefore.After(d) {
			return false
		}
	}
	return true
}

// Len ... Implement sort.Interface
func (evs Events) Len() int {
	return len(evs)
}

// Less ... Implement sort.Interface
func (evs Events) Less(i, j int) bool {
	return evs[i].NotBefore.Unix() < evs[j].NotBefore.Unix()
}

// Swap ... Implement sort.Interface
func (evs Events) Swap(i, j int) {
	evs[i], evs[j] = evs[j], evs[i]
}

// CreateMessage ... Information for displaying to Mackerel
func (evs Events) String() string {
	ev := evs.GetCloseEvent()

	if evs.Len() <= 1 {
		return ev.String()
	}

	return fmt.Sprintf(
		"Instances: %v, %v",
		evs.Len(),
		ev.String(),
	)
}
