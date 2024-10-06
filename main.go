package main

import (
	"fmt"
	"sync"

	"github.com/joetifa2003/conc-map/mapcustom"
)

func main() {
	x := mapcustom.New[string, int]()

	var wg sync.WaitGroup
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			x.Set("a", 1)
			x.Set("b", 2)
			x.Set("c", 3)
		}()
	}
	wg.Wait()

	for k, v := range x.Iter() {
		fmt.Println(k, v)
	}
}
