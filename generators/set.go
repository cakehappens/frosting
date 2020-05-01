package generators

//go:generate genny -pkg frosting -in set.go -out ../ingredient_set.go gen "Item=Ingredient"

import (
	"sync"
)

// ItemQueue the queue of Items
type ItemSet struct {
	items map[*Item]bool
	lock  sync.RWMutex
}

// New creates a new ItemQueue
func NewItemSet(items ...*Item) *ItemSet {
	set := &ItemSet{
		items: make(map[*Item]bool),
	}

	for _, item := range items {
		set.items[item] = true
	}

	return set
}

// Enqueue adds an Item to the end of the queue
func (s *ItemSet) Add(t *Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items[t] = true
}

// Dequeue removes an Item from the start of the queue
func (s *ItemSet) Remove(t *Item) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.items, t)
}

func (s *ItemSet) Contains(t *Item) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.items[t]
}

// Front returns the item next in the queue, without removing it
func (s *ItemSet) Cardinality() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}

func (s *ItemSet) Items() []*Item {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var itemList []*Item
	for key, _ := range s.items {
		itemList = append(itemList, key)
	}

	return itemList
}

// Set Difference: The relative complement or set difference of sets A and B,
// denoted A – B, is the set of all elements in A that are not in B.
// In set-builder notation, A – B = {x ∈ U : x ∈ A and x ∉ B}= A ∩ B'.
func (s *ItemSet) Difference(other *ItemSet) *ItemSet {
	s.lock.RLock()
	defer s.lock.RUnlock()

	diff := NewItemSet()

	for key, _ := range s.items {
		if !other.items[key] {
			diff.items[key] = true
		}
	}

	return diff
}

func (s *ItemSet) Clone() *ItemSet {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return NewItemSet(s.Items()...)
}
