package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

const MAX_CONCURRENCY = 5
const MAX_PAGES = 25

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("invalid starting url %v\n", err)
		os.Exit(1)
	}

	cfg := config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, MAX_CONCURRENCY),
		wg:                 &sync.WaitGroup{},
		maxPages:           MAX_PAGES,
	}

	fmt.Printf("starting crawl of: %v\n", baseURL)
	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL.String())
	cfg.wg.Wait()

	for k, v := range cfg.pages {
		fmt.Println(k, v)
	}

}
