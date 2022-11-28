package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var i int32 = 0

func main() {
	// parallel(plusplus1)
	// fmt.Println(i)
	parallel(plusplus2)
	fmt.Println(i)
}

func plusplus1() int32 {
	i++
	j := i
	return j
}

func plusplus2() int32 {
	j := atomic.AddInt32(&i, 1)
	return j
}

func parallel(f func() int32) {
	var wg sync.WaitGroup
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}
	wg.Wait()
}
