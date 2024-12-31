package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func normalizeURL(inputUrl string) (string, error) {
	URL, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	outUrl, err := url.JoinPath(URL.Hostname(), URL.Path)
	if err != nil {
		return "", err
	}

	outUrl = strings.TrimRight(outUrl, "/")

	outUrl = strings.ToLower(outUrl)

	return outUrl, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)

	document, err := html.Parse(htmlReader)
	if err != nil {
		return []string{}, err
	}

	var links []string
	for node := range document.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, a := range node.Attr {
				if a.Key == "href" {
					var newLink string
					if a.Val[0] == '/' {
						newLink, err = url.JoinPath(rawBaseURL, a.Val)
						if err != nil {
							return []string{}, err
						}
					} else {
						newLink = a.Val
					}
					links = append(links, newLink)
					break
				}
			}

		}
	}

	return links, nil
}
