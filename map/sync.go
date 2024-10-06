package concmap

import "sync"

type syncMap[K comparable, V any] struct {
	m *sync.Map
}

func newSyncMap[K comparable, V any]() syncMap[K, V] {
	return syncMap[K, V]{
		m: &sync.Map{},
	}
}

func (m *syncMap[K, V]) set(k K, v V) {
	m.m.Store(k, v)
}

func (m *syncMap[K, V]) get(k K) (V, bool) {
	var zero V

	res, ok := m.m.Load(k)
	if res == nil {
		return zero, false
	}

	return res.(V), ok
}

func (m *syncMap[K, V]) delete(k K) {
	m.m.Delete(k)
}

func (m *syncMap[K, V]) forEach(f func(k K, v V) bool) {
	m.m.Range(func(k, v any) bool {
		if !f(k.(K), v.(V)) {
			return false
		}

		return true
	})
}
