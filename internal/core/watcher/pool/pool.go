package pool

import (
	"sync"
)

type Pool struct {
	mu			sync.RWMutex
	messages    chan Message
	subscribers []Subscriber
}

func NewPool() *Pool {
	return &Pool{
		messages:    make(chan Message, 128),
		subscribers: make([]Subscriber, 0),
	}
}

func (p *Pool) AddMessage(message Message) {
	select {
	case p.messages <- message:
		p.notifySubscribers(message)
	default:
	}
}

func (p *Pool) Subscribe(callback Subscriber) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers = append(p.subscribers, callback)
}

func (p *Pool) notifySubscribers(message Message) {
	go func() {
		p.mu.RLock()
		subscribers := make([]Subscriber, len(p.subscribers))
		copy(subscribers, p.subscribers)
		p.mu.RUnlock()
		
		for _, subscriber := range subscribers {
			subscriber(message.Level, message.Message)
		}
	}()
}

func (p *Pool) GetSubscriberCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.subscribers)
}
