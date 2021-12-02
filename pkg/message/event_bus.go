package message

import (
	"context"
)

type Event uint32

const (
	NewMessageEvent Event = iota
	InvalidEvent
)

func (e Event) String() string {
	switch e {
	case NewMessageEvent:
		return "NewMessage"
	default:
		return "Invalid"
	}
}

type EventHnd func(ctx context.Context, res *TicketMessage)

type EventBus interface {
	Subscribe(Event, EventHnd)
	Publish(Event, context.Context, *TicketMessage)
}
