package list

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

const GOROUTINES = 100

const OPS = 1000

func TestList(t *testing.T) {
	ll := NewLinkedList[string]()

	var wg sync.WaitGroup
	for i := range GOROUTINES {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			for j := range OPS {
				ll.Insert(fmt.Sprintf("%d-%d", i, j))
			}
		}(i)
	}

	wg.Wait()

	for i := range GOROUTINES {
		for j := range OPS {
			if !ll.Contains(fmt.Sprintf("%d-%d", i, j)) {
				t.Errorf("Expected %s to be in the list", fmt.Sprintf("%d-%d", i, j))
			}
		}
	}
}

func BenchmarkLinkedList(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func(b *testing.B, numGoroutines int)
	}{
		{"LockFree", benchmarkLockFree},
		{"Locking", benchmarkLocking},
	}

	goroutines := []int{1, 2, 4, 8, 16, 32}

	for _, bm := range benchmarks {
		for _, numGoroutines := range goroutines {
			name := fmt.Sprintf("%s-%d", bm.name, numGoroutines)
			b.Run(name, func(b *testing.B) {
				bm.fn(b, numGoroutines)
			})
		}
	}
}

func benchmarkLockFree(b *testing.B, numGoroutines int) {
	list := NewLinkedList[int]()
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			r := rand.New(rand.NewSource(int64(i)))
			for j := 0; j < b.N/numGoroutines; j++ {
				op := r.Intn(3)
				value := r.Intn(1000)
				switch op {
				case 0:
					list.Insert(value)
				case 1:
					list.Remove(value)
				case 2:
					list.Contains(value)
				}
			}
		}()
	}
	wg.Wait()
}

func benchmarkLocking(b *testing.B, numGoroutines int) {
	list := NewLockingLinkedList[int]()
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			r := rand.New(rand.NewSource(int64(i)))
			for j := 0; j < b.N/numGoroutines; j++ {
				op := r.Intn(3)
				value := r.Intn(1000)
				switch op {
				case 0:
					list.Insert(value)
				case 1:
					list.Remove(value)
				case 2:
					list.Contains(value)
				}
			}
		}()
	}
	wg.Wait()
}
