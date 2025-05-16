package search

import (
	"sync"

	"github.com/google/uuid"
)

type Blog struct {
	ID           uuid.UUID
	Title        string
	Brief        string
	ThumbnailUrl string
	Views        int32
}

type Book struct {
	ID            uuid.UUID
	Name          string
	CoverImageUrl string
}

type AllowedTypes interface {
	Book | Blog
}

type safeSet[T AllowedTypes] struct {
	mutex sync.RWMutex
	items map[T]struct{}
}

func NewSafeSet[T AllowedTypes]() *safeSet[T] {
	return &safeSet[T]{
		items: make(map[T]struct{}),
	}
}

func (ss *safeSet[T]) Add(value T) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	if _, ok := ss.items[value]; !ok {
		ss.items[value] = struct{}{}
	}
}

func (ss *safeSet[T]) Has(value T) bool {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	_, ok := ss.items[value]
	return ok
}

func (ss *safeSet[T]) Keys() []*T {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	var keys []*T
	for key := range ss.items {
		keys = append(keys, &key)
	}

	if len(keys) > 0 {
		return keys
	}

	return make([]*T, 0)
}

func (ss *safeSet[T]) Remove(value T) bool {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	if exists := ss.Has(value); exists {
		delete(ss.items, value)
		return true
	}

	return false
}
