package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", errors.New("unable to fetch site")
	}

	header := resp.Header.Get("content-type")
	if !strings.Contains(header, "text/html") {
		return "", fmt.Errorf("unable parse site data: %v", header)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(dat), nil

}
