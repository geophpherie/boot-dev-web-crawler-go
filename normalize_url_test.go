package main

import (
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme and trail slash",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme secure",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove scheme secure and trail slash",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "ALL CAPS URL",
			inputURL: "HTTPS://BLOG.BOOT.DEV/PATH/",
			expected: "blog.boot.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%v' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - '%v' FAIL: expected URL: %v, GOT: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "absolute and relative URLs more nested",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<h1><a href="/path/one">
			<span>Boot.dev</span>
		</a></h1>
		<h2><a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a></h2>
	</body>
</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%v' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if isEqual := reflect.DeepEqual(actual, tc.expected); !isEqual {
				t.Errorf("Test %v - '%v' FAIL: EXPECTED: %v, GOT: %v", i, tc.name, tc.expected, actual)
			}
		})
	}

}
