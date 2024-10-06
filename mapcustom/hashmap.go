package mapcustom

import (
	"iter"
	"sync/atomic"

	"github.com/benbjohnson/immutable"
	"github.com/dolthub/maphash"

	"github.com/joetifa2003/conc-map/cow"
	"github.com/joetifa2003/conc-map/list"
)

type LFMap[K comparable, V comparable] struct {
	buckets *list.LinkedList[bucket[K, V]]
	h       maphash.Hasher[K]
	total   *atomic.Int32
}

type bucket[K comparable, V comparable] struct {
	*cow.Cow[immutable.List[pair[K, V]]]
}

func newBucket[K comparable, V comparable]() bucket[K, V] {
	return bucket[K, V]{Cow: cow.NewCow(immutable.NewList[pair[K, V]]())}
}

type pair[K comparable, V comparable] struct {
	key   K
	value V
}

func New[K comparable, V comparable]() *LFMap[K, V] {
	buckets := list.NewLinkedList[bucket[K, V]]()
	buckets.Insert(newBucket[K, V]())

	total := atomic.Int32{}
	total.Store(1)

	return &LFMap[K, V]{
		buckets: buckets,
		total:   &total,
		h:       maphash.NewHasher[K](),
	}
}

func (m *LFMap[K, V]) Set(key K, value V) {
	bucketIndex := m.getBucketIndex(key)
	bucket := m.buckets.Get(bucketIndex)
	bucket.Tx(func(l *immutable.List[pair[K, V]]) *immutable.List[pair[K, V]] {
		itr := l.Iterator()
		for !itr.Done() {
			idx, val := itr.Next()
			if val.key == key {
				return l.Set(idx, pair[K, V]{key: key, value: value})
			}
		}

		return l.Append(pair[K, V]{key: key, value: value})
	})
}

func (m *LFMap[K, V]) Get(key K) (V, bool) {
	bucketIndex := m.getBucketIndex(key)
	bucket := m.buckets.Get(bucketIndex)

	var zero V

	itr := bucket.Get().Iterator()
	for !itr.Done() {
		_, val := itr.Next()
		if val.key == key {
			return val.value, true
		}
	}

	return zero, false
}

func (m *LFMap[K, V]) getBucketIndex(key K) int {
	return int(m.h.Hash(key) % uint64(m.buckets.Len()))
}

func (m *LFMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for b := range m.buckets.Iter() {

			itr := b.Get().Iterator()
			for !itr.Done() {
				_, val := itr.Next()
				if !yield(val.key, val.value) {
					return
				}
			}

		}
	}
}
