package main

import (
	"fmt"

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
	ConcurrentMutex("http://golang.org/", fetcher, makeState())

	// fmt.Printf("=== ConcurrentChannel ===\n")
	// ConcurrentChannel("http://golang.org/", fetcher)
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
