package list

import (
	"iter"
	"sync"
)

type LockingNode[T comparable] struct {
	value T
	next  *LockingNode[T]
}

type LockingLinkedList[T comparable] struct {
	head *LockingNode[T]
	mu   sync.RWMutex
}

func NewLockingLinkedList[T comparable]() *LockingLinkedList[T] {
	return &LockingLinkedList[T]{}
}

func (ll *LockingLinkedList[T]) Insert(value T) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	newNode := &LockingNode[T]{value: value, next: ll.head}
	ll.head = newNode
}

func (ll *LockingLinkedList[T]) Remove(value T) bool {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	if ll.head == nil {
		return false
	}
	if ll.head.value == value {
		ll.head = ll.head.next
		return true
	}
	curr := ll.head
	for curr.next != nil {
		if curr.next.value == value {
			curr.next = curr.next.next
			return true
		}
		curr = curr.next
	}
	return false
}

func (ll *LockingLinkedList[T]) Contains(value T) bool {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	curr := ll.head
	for curr != nil {
		if curr.value == value {
			return true
		}
		curr = curr.next
	}
	return false
}

func (ll *LockingLinkedList[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		ll.mu.RLock()
		defer ll.mu.RUnlock()
		current := ll.head
		for current != nil {
			if !yield(current.value) {
				return
			}
			current = current.next
		}
	}
}
