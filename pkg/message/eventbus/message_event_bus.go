package eventbus

import (
	"context"
	"sync"

	"github.com/wascript3r/autonuoma/pkg/message"
	"github.com/wascript3r/cryptopay/pkg/logger"
	"github.com/wascript3r/gopool"
)

type EventBus struct {
	pool *gopool.Pool
	log  logger.Usecase

	mx       *sync.RWMutex
	handlers map[message.Event][]message.EventHnd
}

func New(pool *gopool.Pool, log logger.Usecase) *EventBus {
	return &EventBus{
		pool: pool,
		log:  log,

		mx:       &sync.RWMutex{},
		handlers: make(map[message.Event][]message.EventHnd),
	}
}

func (e *EventBus) Subscribe(ev message.Event, hnd message.EventHnd) {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.handlers[ev] = append(e.handlers[ev], hnd)
}

func (e *EventBus) Publish(ev message.Event, ctx context.Context, tm *message.TicketMessage) {
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
			h(ctx, tm)
			wg.Done()
		})
		if err != nil {
			e.log.Error("Cannot publish message %s event because of pool schedule error: %s", ev, err)
			wg.Done()
		}
	}

	wg.Wait()
}
