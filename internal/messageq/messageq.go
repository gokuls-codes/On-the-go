package messageq

import (
	"fmt"
	"sync"
)

type MessageQ struct {
	Messages map[string]chan string
	mu       sync.Mutex
}

func NewMessageQ() *MessageQ {
	return &MessageQ{
		Messages: make(map[string]chan string),
	}
}

func (mq *MessageQ) CreateChannel(id string) chan string {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if ch, exists := mq.Messages[id]; exists {
		return ch
	}

	ch := make(chan string)
	mq.Messages[id] = ch
	return ch
}

func (mq *MessageQ) GetChannel(id string) (chan string, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if ch, exists := mq.Messages[id]; exists {
		return ch, nil
	}

	return nil, fmt.Errorf("channel not found")
}

func (mq *MessageQ) RemoveChannel(id string) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if ch, exists := mq.Messages[id]; exists {
		close(ch)
		delete(mq.Messages, id)
	}
}
