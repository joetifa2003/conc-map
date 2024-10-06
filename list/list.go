package list

import (
	"iter"
	"sync/atomic"
)

// Node represents a node in the linked list
type Node[T comparable] struct {
	value T
	next  atomic.Pointer[Node[T]]
}

// LinkedList represents the linked list
type LinkedList[T comparable] struct {
	head atomic.Pointer[Node[T]]
	len  atomic.Int32
}

// NewLinkedList creates a new empty linked list
func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Insert adds a new value to the list, reusing removed nodes if available
func (ll *LinkedList[T]) Insert(value T) {
	newNode := &Node[T]{value: value}
	ll.len.Add(1)

	for {
		newNode.next.Store(ll.head.Load())
		if ll.head.CompareAndSwap(newNode.next.Load(), newNode) {
			return
		}
	}
}

// Contains checks if the given value exists in the list (not marked as removed)
func (ll *LinkedList[T]) Contains(value T) bool {
	curr := ll.head.Load()
	for curr != nil {
		if curr.value == value {
			return true
		}
		curr = curr.next.Load()
	}
	return false
}

func (ll *LinkedList[T]) Len() int {
	return int(ll.len.Load())
}

// Iterate traverses the list and calls the provided callback function for each non-removed element
func (ll *LinkedList[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := ll.head.Load()
		for current != nil {
			if !yield(current.value) {
				return
			}
			current = current.next.Load()
		}
	}
}

// Iterate traverses the list and calls the provided callback function for each non-removed element
func (ll *LinkedList[T]) IterPtr() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		current := ll.head.Load()
		for current != nil {
			if !yield(&current.value) {
				return
			}
			current = current.next.Load()
		}
	}
}

func (ll *LinkedList[T]) Get(index int) T {
	if index < 0 || index >= ll.Len() {
		panic("index out of range")
	}

	current := ll.head.Load()
	currentIndex := 0

	for current != nil {
		if currentIndex == index {
			return current.value
		}
		current = current.next.Load()
		currentIndex++
	}

	panic("index out of range")
}

func (ll *LinkedList[T]) Clone() *LinkedList[T] {
	newList := NewLinkedList[T]()

	for v := range ll.Iter() {
		newList.Insert(v)
	}

	return newList
}
