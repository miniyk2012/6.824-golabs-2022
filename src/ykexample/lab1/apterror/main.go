package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go fmt.Println(<-c) // 会dead lock, <-c作用域是跟着外层代码的, 而不是新的协程的
	c <- 5
	time.Sleep(time.Second)
}
