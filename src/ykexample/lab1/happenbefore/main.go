// https://studygolang.com/articles/14129

package main

import (
	"fmt"
	"time"
)

func main() {
	happensBeforeMulti(10)
	fmt.Println("---")
	c1() // hello, world
	fmt.Println("---")
	c2() // 有可能为空
	fmt.Println("---")
	c3() // hello, world
}

// Sample Routine 1
func happensBeforeMulti(i int) {
	i += 2      // E1
	go func() { // G1 goroutine create
		fmt.Println(i) // E2
	}() // G2 goroutine destryo
	time.Sleep(time.Millisecond * 100)
}

// receive早于send完成
func c1() {
	var c = make(chan int)
	var a string
	f := func() {
		a = "hello, world"
		<-c // receive
	}
	go f()
	c <- 0
	println(a) // send完成
}

// 第k个receive早于第k+10个send完成
func c2() {
	// Channel routine 2
	var c = make(chan int, 10)
	var a string

	f := func() {
		a = "hello, world"
		<-c // E1 receive
	}

	go f()
	c <- 0     // E2 bufferd channel send完成
	println(a) // E1不一定 happens Before E2
	time.Sleep(time.Millisecond * 100)
}

// send早于receive完成
func c3() {
	// Channel routine 3
	var c = make(chan int, 10)
	var a string

	f := func() {
		a = "hello, world"
		c <- 0 // send
	}

	go f()
	<-c // receive完成
	println(a)
}
