package concmap

import (
	"sync"
	"time"

	"github.com/dolthub/maphash"
)

type Option[K comparable, V any] func(m *Map[K, V])

func WithShardCount[K comparable, V any](n int) Option[K, V] {
	return func(m *Map[K, V]) {
		m.shardCount = n
	}
}

type atom[V any] struct {
	v V
	l *sync.RWMutex
}

type Map[K comparable, V any] struct {
	shards     []syncMap[K, V]
	shardCount int
	mh         maphash.Hasher[K]
}

func New[K comparable, V any](opts ...Option[K, V]) *Map[K, V] {
	m := &Map[K, V]{
		shardCount: 1,
		mh:         maphash.NewHasher[K](),
	}

	for _, o := range opts {
		o(m)
	}

	m.shards = make([]syncMap[K, V], m.shardCount)

	for i := range m.shards {
		m.shards[i] = newSyncMap[K, V]()
	}

	return m
}

type setOptions struct {
	ttl time.Duration
}

type SetOption func(s *setOptions)

func WithTTL(ttl time.Duration) SetOption {
	return func(s *setOptions) {
		s.ttl = ttl
	}
}

func (m *Map[K, V]) Set(k K, v V, opts ...SetOption) {
	s := &setOptions{}

	for _, o := range opts {
		o(s)
	}

	shard := m.getShard(k)

	shard.set(k, v)

	if s.ttl != 0 {
		time.AfterFunc(s.ttl, func() {
			m.Delete(k)
		})
	}
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	shard := m.getShard(k)
	return shard.get(k)
}

func (m *Map[K, V]) Delete(k K) {
	shard := m.getShard(k)
	shard.delete(k)
}

func (m *Map[K, V]) Len() int {
	res := 0

	// for _, s := range m.shards {
	// 	res += len(s.m)
	// }

	return res
}

func (m *Map[K, V]) ForEach(f func(k K, v V) bool) {
	for _, s := range m.shards {
		s.forEach(func(k K, v V) bool {
			if !f(k, v) {
				return false
			}

			return true
		})
	}
}

func (m *Map[K, V]) getShard(k K) syncMap[K, V] {
	shardIdx := m.mh.Hash(k) % uint64(len(m.shards))
	return m.shards[shardIdx]
}
