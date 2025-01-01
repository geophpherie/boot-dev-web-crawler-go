package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]

	fmt.Printf("starting crawl of: %v\n", baseURL)

	pages := make(map[string]int)
	pages, err := crawlPage(baseURL, baseURL, pages)
	if err != nil {
		fmt.Printf("crawl failed %v\n", err)
		os.Exit(1)
	}

	for k, v := range pages {
		fmt.Println(k, v)
	}

}
