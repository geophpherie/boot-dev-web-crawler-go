package main

import (
	"cmp"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

var MAX_CONCURRENCY int
var MAX_PAGES int

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) == 1 {
		MAX_CONCURRENCY = 5
		MAX_PAGES = 20
	} else {
		MAX_CONCURRENCY, _ = strconv.Atoi(args[1])
		MAX_PAGES, _ = strconv.Atoi(args[2])
	}

	fmt.Println(MAX_CONCURRENCY, MAX_PAGES)

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

	fmt.Printf(`
=============================
  REPORT for %v
=============================
`, cfg.baseURL.String())

	reportLinks := []reportLink{}

	for k, v := range cfg.pages {
		reportLinks = append(reportLinks, reportLink{count: v, url: k})
	}

	slices.SortFunc(reportLinks, func(a, b reportLink) int {
		// negative when a < b
		// positive when a > b
		// 0 when a ==b
		if n := cmp.Compare(b.count, a.count); n != 0 {
			return n
		}

		return strings.Compare(a.url, b.url)
	})

	for _, l := range reportLinks {
		fmt.Printf("Found %v internal links to %v\n", l.count, l.url)
	}
}

type reportLink struct {
	count int
	url   string
}
