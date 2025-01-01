package main

import (
	"reflect"
	"testing"
)

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
