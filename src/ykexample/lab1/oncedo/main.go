package main

import (
	"fmt"
	"sync"
)

var once sync.Once

func main() {
	once.Do(func() {
		fmt.Println("ha")
	})
	once.Do(func() {
		fmt.Println("never run")
	})
}
