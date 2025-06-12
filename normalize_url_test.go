package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normalize_url",
			input:    "https://www.google.com",
			expected: "https://www.google.com",
		},
		{
			name:     "normalize_url_with_trailing_slash",
			input:    "https://www.google.com/",
			expected: "https://www.google.com",
		},
		{
			name:     "normalize_url_with_query_params",
			input:    "https://www.google.com?q=test",
			expected: "https://www.google.com?q=test",
		},
		{
			name:     "normalize_url_with_fragment",
			input:    "https://www.google.com#test",
			expected: "https://www.google.com",
		},
		{
			name:     "normalize_url_with_path_traversal",
			input:    "https://WWW.Example.com:443/../a/./b/index.html?b=2&a=1#section",
			expected: "https://www.example.com/a/b/index.html?a=1&b=2",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.input)
			if err != nil {
				t.Errorf("test %d: expected no error, got %v", i, err)
			}
			if actual != tc.expected {
				t.Errorf("test %d: expected %s, got %s", i, tc.expected, actual)
			}
		})
	}
}
