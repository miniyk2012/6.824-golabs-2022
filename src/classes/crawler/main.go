package main

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/elliotchance/pie/v2"
)

//
// main
//

func main() {
	fmt.Printf("=== Serial===\n")
	fetched := make(map[string]bool)
	Serial("http://golang.org/", fetcher, fetched)
	spew.Printf("total urls=%+v\n", pie.Keys(fetched))

	fmt.Printf("=== ConcurrentMutex ===\n")
	state := makeState()
	ConcurrentMutex("http://golang.org/", fetcher, state)
	spew.Printf("total urls=%+v\n", pie.Keys(state.fetched))

	fmt.Printf("=== ConcurrentChannel ===\n")
	ConcurrentChannel("http://golang.org/", fetcher)
}

func ConcurrentChannel(url string, fetcher Fetcher) {

	c := make(chan []string)
	go func() {
		c <- []string{"http://golang.org/"}
	}()
	coordinator(c, fetcher)
}

func coordinator(c chan []string, fetcher Fetcher) {
	fetched := make(map[string]bool)
	n := 1 // 记录启动了多少个worker
	for urls := range c {
		for _, url := range urls {
			if !fetched[url] {
				fetched[url] = true
				n++
				go worker(url, c, fetcher)
			}
		}
		n--
		// 当所有worker都返回
		if n == 0 {
			break
		}
	}
	// spew.Printf("total urls=%+v\n", pie.Keys(fetched))
}

func worker(url string, c chan<- []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	if err != nil {
		c <- []string{}
	} else {
		c <- urls
	}
}

//
// Several solutions to the crawler exercise from the Go tutorial
// https://tour.golang.org/concurrency/10
//

//
// Serial crawler
//
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("fetch failed: %v\n", err)
		return
	}
	fetched[url] = true
	for _, url := range urls {
		Serial(url, fetcher, fetched)
	}
}

//
// Concurrent crawler with shared state and Mutex
//

type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func makeState() *fetchState {
	state := &fetchState{}
	state.fetched = make(map[string]bool)
	return state
}

func ConcurrentMutex(url string, fetcher Fetcher, state *fetchState) {
	state.mu.Lock()
	if state.fetched[url] {
		state.mu.Unlock()
		return
	}
	urls, err := fetcher.Fetch(url)
	if err != nil {
		state.mu.Unlock()
		fmt.Printf("fetch failed: %v\n", err)
		return
	}
	state.fetched[url] = true
	state.mu.Unlock()
	var done sync.WaitGroup
	for _, aurl := range urls {
		done.Add(1)
		go func(url string) {
			ConcurrentMutex(url, fetcher, state)
			done.Done()
		}(aurl)
	}
	done.Wait()
	return
}
