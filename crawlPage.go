package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	// aquire mutex
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	// if page exists in map, increment
	if _, ok := cfg.pages[normalizedURL]; ok {
		// if not init it and set to 1
		cfg.pages[normalizedURL] += 1
		return false
	} else {
		cfg.pages[normalizedURL] = 1
		return true
	}
}

func (cfg *config) crawlPage(rawCurrentURL string) error {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return nil
	}
	cfg.mu.Unlock()

	fmt.Printf("Crawling %v\n", rawCurrentURL)
	// make sure on same domain
	currentUrl, err := url.Parse(rawCurrentURL)
	if err != nil {
		return err
	}

	if cfg.baseURL.Hostname() != currentUrl.Host {
		return err
	}

	// normalize current url
	normalizedCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return err
	}

	isFirstVisit := cfg.addPageVisit(normalizedCurrent)
	if !isFirstVisit {
		return nil
	}

	// get the htm for the url
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return err
	}

	// get all urls in the body
	urls, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		return err
	}

	// crawl any urls
	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}

	return nil
}
