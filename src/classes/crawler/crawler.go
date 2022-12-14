package main

import (
	"fmt"
)

/*
func ConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	f.mu.Lock()
	already := f.fetched[url]
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, f)
		}(u)
	}
	done.Wait()
	return
}
*/

//
// Concurrent crawler with channels
//

// func worker(url string, ch chan []string, fetcher Fetcher) {
// 	urls, err := fetcher.Fetch(url)
// 	if err != nil {
// 		ch <- []string{}
// 	} else {
// 		ch <- urls
// 	}
// }

// func coordinator(ch chan []string, fetcher Fetcher) {
// 	n := 1 // 记录启动了多少个worker
// 	fetched := make(map[string]bool)
// 	for urls := range ch {
// 		for _, u := range urls {
// 			if fetched[u] == false {
// 				fetched[u] = true
// 				n += 1
// 				go worker(u, ch, fetcher)
// 			}
// 		}
// 		n -= 1
// 		if n == 0 { // 当所有worker都返回
// 			break
// 		}
// 	}
// }

// func ConcurrentChannel(url string, fetcher Fetcher) {
// 	ch := make(chan []string)
// 	go func() {
// 		ch <- []string{url}
// 	}()
// 	coordinator(ch, fetcher)
// }

//
// Fetcher
//

type Fetcher interface {
	// Fetch returns a slice of URLs found on the page.
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		fmt.Printf("found:   %s\n", url)
		return res.urls, nil
	}
	fmt.Printf("missing: %s\n", url)
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
