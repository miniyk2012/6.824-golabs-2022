package main

import (
	"fmt"
	"time"
)

func main() {
	var t *time.Time
	fmt.Println(t)
	a := time.Now()
	t = &a
	fmt.Println(t)
}
