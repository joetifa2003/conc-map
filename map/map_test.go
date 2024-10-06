package concmap

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
)

type testStruct struct {
	A string
	B int
	F float32
}

const GOROUTINES = 100

const ITERATIONS = 500

const PERGOROUTINE = ITERATIONS / GOROUTINES

func variableShards(shardCounts []int, b *testing.B, f func(b *testing.B, m *Map[string, testStruct])) {
	for _, s := range shardCounts {
		b.Run(fmt.Sprintf("concmap_shards_%d", s), func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				m := New[string, testStruct](
					WithShardCount[string, testStruct](s),
				)

				f(b, m)
			}
		})
	}
}

func BenchmarkVsSync(b *testing.B) {
	variableShards([]int{1, 10, 50}, b, func(b *testing.B, m *Map[string, testStruct]) {
		rand := rand.New(rand.NewPCG(0, 0))

		var wg sync.WaitGroup
		for _ = range GOROUTINES {
			wg.Add(1)

			go func() {
				defer wg.Done()
				for range PERGOROUTINE {
					s := fmt.Sprint(rand.Int())
					m.Set(s, testStruct{A: s, B: rand.Int(), F: float32(rand.Float32())})

					m.ForEach(func(k string, v testStruct) bool {
						return true
					})
				}
			}()
		}

		wg.Wait()
	})

	b.Run("sync.Map", func(b *testing.B) {
		for range b.N {
			m := sync.Map{}

			rand := rand.New(rand.NewPCG(0, 0))
			var wg sync.WaitGroup

			for _ = range GOROUTINES {
				wg.Add(1)

				go func() {
					defer wg.Done()
					for range PERGOROUTINE {
						s := fmt.Sprint(rand.Int())
						m.Store(s, testStruct{A: s, B: rand.Int(), F: float32(rand.Float32())})

						m.Range(func(k any, v any) bool {
							return true
						})
					}
				}()

			}

			wg.Wait()
		}
	})
}
