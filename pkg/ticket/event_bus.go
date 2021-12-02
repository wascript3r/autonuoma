package ticket

import (
	"context"
)

type Event uint32

const (
	NewTicketEvent Event = iota
	AcceptedTicketEvent
	EndedTicketEvent
	InvalidEvent
)

func (e Event) String() string {
	switch e {
	case NewTicketEvent:
		return "NewTicket"
	case AcceptedTicketEvent:
		return "AcceptedTicket"
	case EndedTicketEvent:
		return "EndedTicket"
	default:
		return "Invalid"
	}
}

type EventHnd func(ctx context.Context, ticketID int)

type EventBus interface {
	Subscribe(Event, EventHnd)
	Publish(Event, context.Context, int)
}
