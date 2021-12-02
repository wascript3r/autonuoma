package eventbus

import (
	"context"
	"sync"

	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/autonuoma/pkg/ticket"
	"github.com/wascript3r/cryptopay/pkg/logger"
	"github.com/wascript3r/gopool"
)

type EventBus struct {
	pool *gopool.Pool
	log  logger.Usecase

	mx       *sync.RWMutex
	handlers map[ticket.Event][]ticket.EventHnd
}

func New(pool *gopool.Pool, log logger.Usecase) *EventBus {
	return &EventBus{
		pool: pool,
		log:  log,

		mx:       &sync.RWMutex{},
		handlers: make(map[ticket.Event][]ticket.EventHnd),
	}
}

func (e *EventBus) Subscribe(ev ticket.Event, hnd ticket.EventHnd) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.handlers[ev] = append(e.handlers[ev], hnd)
}

func (e *EventBus) Publish(ev ticket.Event, ctx context.Context, ticketID int, tm *message.TicketMessage) {
	e.mx.RLock()
	defer e.mx.RUnlock()

	hnds := e.handlers[ev]
	count := len(hnds)
	if count == 0 {
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(count)

	for _, h := range hnds {
		h := h
		err := e.pool.Schedule(func() {
			h(ctx, ticketID, tm)
			wg.Done()
		})
		if err != nil {
			e.log.Error("Cannot publish ticket %s event because of pool schedule error: %s", ev, err)
			wg.Done()
		}
	}

	wg.Wait()
}
