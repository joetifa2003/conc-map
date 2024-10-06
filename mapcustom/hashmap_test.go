package mapcustom

import (
	"math/rand/v2"
	"sync"
	"testing"

	"github.com/Snawoot/lfmap"
)

// goos: linux
// goarch: amd64
// pkg: github.com/joetifa2003/conc-map/mapcustom
// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLFMap-16            865           1445544 ns/op         3213814 B/op      27781 allocs/op
// BenchmarkSync-16            1887            653815 ns/op          175533 B/op       4147 allocs/op
const GOROUTINES = 100

const ITERATIONS = 10

func benchPar[T any](b *testing.B, factory func() T, f func(b *testing.B, m T)) {
	for range b.N {
		m := factory()
		var wg sync.WaitGroup
		for range GOROUTINES {
			wg.Add(1)
			go func() {
				defer wg.Done()
				f(b, m)
			}()
		}
		wg.Wait()
	}
}

func BenchmarkLFMap(b *testing.B) {
	benchPar(b,
		func() *lfmap.Map[int, int] { return lfmap.New[int, int]() },
		func(b *testing.B, m *lfmap.Map[int, int]) {
			for i := 0; i < ITERATIONS; i++ {
				m.Set(rand.Int(), i)
			}
		},
	)
}

func BenchmarkSync(b *testing.B) {
	benchPar(b,
		func() *sync.Map { return &sync.Map{} },
		func(b *testing.B, m *sync.Map) {
			for i := 0; i < ITERATIONS; i++ {
				m.Store(rand.Int(), i)
			}
		},
	)
}

func BenchmarkConcMap(b *testing.B) {
	benchPar(b,
		func() *LFMap[int, int] { return New[int, int]() },
		func(b *testing.B, m *LFMap[int, int]) {
			for i := 0; i < ITERATIONS; i++ {
				m.Set(rand.Int(), i)
			}
		},
	)
}
