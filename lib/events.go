package checkawsec2mainte

import (
	"sort"
	"strings"
	"time"

	"github.com/ntrv/check-aws-ec2-mainte/lib/ec2api"
	"github.com/ntrv/check-aws-ec2-mainte/lib/metadata"
)

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

// UpdateStates ... Descriptionに含まれている文字列からStateを設定
func (evs Events) UpdateStates() {
	for i, ev := range evs {
		switch {
		case strings.Contains(ev.Description, "[Completed]"):
			evs[i].State = StateCompleted
		case strings.Contains(ev.Description, "[Canceled]"):
			evs[i].State = StateCanceled
		default:
			evs[i].State = StateActive
		}
	}
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

func (evs Events) SetMetadataEvents(events metadata.Events) {
	for _, event := range events {
		evs = append(evs, Event{
			Code:        event.Code,
			InstanceId:  event.InstanceId,
			NotBefore:   time.Time(event.NotBefore),
			NotAfter:    time.Time(event.NotAfter),
			Description: event.Description,
			State:       EventState(event.State),
		})
	}
}

func (evs Events) SetEC2APIEvents(events ec2api.Events) {
	for _, event := range events {
		evs = append(evs, Event{
			Code:        event.Code,
			InstanceId:  event.InstanceId,
			NotBefore:   event.NotBefore,
			NotAfter:    event.NotAfter,
			Description: event.Description,
		})
	}
	evs.UpdateStates()
}
