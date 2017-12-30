package main

import (
	"fmt"
	"sync"
)

var urlState = make(map[string]bool)
var mutex = sync.Mutex{}

type Fetcher interface {
	Fetch(url string, urls chan<- string)
}

func setState(url string) {
 	mutex.Lock()
 	urlState[url] = true
 	mutex.Unlock()
}

func Crawl(url string, depth int, fetcher Fetcher, quit chan bool) {
	if depth <= 0 {
		quit <- true
	}
	
	urls := make(chan string)
	go fetcher.Fetch(url, urls)
	
	for {
		select {
	    case msg := <-urls:
        if !urlState[msg] {
          fmt.Println("received message", msg)
				  setState(msg)
				}
			  go Crawl(msg, depth-1, fetcher, quit)
			case <-quit:
			  fmt.Println("Time to stop")
        return
    }
	}
}

func main() {
	quit := make(chan bool)
	Crawl("http://golang.org/", 4, fetcher, quit)
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string, urls chan<- string) {
		if res, ok := f[url]; ok {
			for _, u := range res.urls {
				urls <- u
			}
		}
}

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
