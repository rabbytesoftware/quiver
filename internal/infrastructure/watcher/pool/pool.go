package pool

import (
	"sync"
)

type Pool struct {
	messages    chan string
	subscribers []func(string)
	mu          sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		messages:    make(chan string, 100),
		subscribers: make([]func(string), 0),
	}
}

func (p *Pool) AddMessage(message string) {
	select {
	case p.messages <- message:
		p.notifySubscribers(message)
	default:
	}
}

func (p *Pool) Subscribe(callback func(string)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subscribers = append(p.subscribers, callback)
}

func (p *Pool) notifySubscribers(message string) {
	p.mu.RLock()
	subscribers := make([]func(string), len(p.subscribers))
	copy(subscribers, p.subscribers)
	p.mu.RUnlock()
	
	for _, subscriber := range subscribers {
		go subscriber(message)
	}
}

func (p *Pool) GetSubscriberCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.subscribers)
}
