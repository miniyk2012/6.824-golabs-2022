package main

import "sync"

var a string
var done bool
var mutex sync.Mutex

func setup() {
	mutex.Lock()
	defer mutex.Unlock()
	done = true // 这2条语句可能重排, 因此加个锁比较安全
	a = "hello world"
}

func main() {
	go setup()

	for {
		mutex.Lock()
		if done {
			mutex.Unlock()
			break
		}
		mutex.Unlock()
	}
	mutex.Lock()
	print(a)
	mutex.Unlock()
}
