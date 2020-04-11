package generators

//go:generate genny -pkg frosting -in queue.go -out ../ingredient_queue.go gen "Item=Ingredient"

import (
	"sync"
)

// ItemQueue the queue of Items
type ItemQueue struct {
	items []*Item
	lock  sync.RWMutex
}

// New creates a new ItemQueue
func NewItemQueue() *ItemQueue {
	return &ItemQueue{
		items: []*Item{},
	}
}

// Enqueue adds an Item to the end of the queue
func (s *ItemQueue) Enqueue(t *Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = append(s.items, t)
}

// Dequeue removes an Item from the start of the queue
func (s *ItemQueue) Dequeue() (*Item, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.items) > 0 {
		item := s.items[0]
		s.items = s.items[1:len(s.items)]
		return item, true
	}

	return nil, false
}

// Front returns the item next in the queue, without removing it
func (s *ItemQueue) Front() *Item {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items[0]
}

// IsEmpty returns true if the queue is empty
func (s *ItemQueue) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Length() == 0
}

// Size returns the number of Items in the queue
func (s *ItemQueue) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
