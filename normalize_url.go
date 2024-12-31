package main

import (
	"net/url"
	"strings"
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
	return outUrl, nil
}
