package main

import (
	"fmt"
	"os"
)

func main() {
	f1()
	fmt.Println("wao")
}

func f1() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("worker[%d] panic, %v", os.Getpid(), r)
			panic(r)
		}
	}()
	f2()
}

func f2() {
	// panic("hello world\n")
	os.Exit(1)
}
