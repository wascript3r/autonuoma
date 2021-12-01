package ticket

import (
	"context"
)

type Event uint32

const (
	NewTicketEvent Event = iota
	InvalidEvent
)

func (e Event) String() string {
	switch e {
	case NewTicketEvent:
		return "NewTicket"
	default:
		return "Invalid"
	}
}

type EventHnd func(ctx context.Context)

type EventBus interface {
	Subscribe(Event, EventHnd)
	Publish(Event, context.Context)
}
