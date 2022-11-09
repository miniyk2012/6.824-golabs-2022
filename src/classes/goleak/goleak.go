// Guest Lecture on Go: Russ Cox: https://www.youtube.com/watch?v=IdCbMO0Ey9I

package main

import (
	"net/http"
	_ "net/http/pprof"
)

var c = make(chan int)

func main() {
	for i := 0; i < 100; i++ {
		go f(0x10 * i)
	}
	http.ListenAndServe("localhost:8080", nil) // 为了不退出main, 并启动pprof的服务
	// http://localhost:8080/debug/pprof/goroutine?debug=1
}

func f(x int) {
	g(x + 1)
}

func g(x int) {
	h(x + 1)
}

func h(x int) {
	c <- 1
	f(x + 1)
}
