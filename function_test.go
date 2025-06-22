package main

import (
	"reflect"
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
		{
			name:     "normalize_url_with_multiple_slashes",
			input:    "https://example.com//path//to//file",
			expected: "https://example.com/path/to/file",
		},
		{
			name:     "normalize_url_with_percent_encoding",
			input:    "https://example.com/path%20with%20spaces",
			expected: "https://example.com/path%20with%20spaces",
		},
		{
			name:     "normalize_url_with_empty_path",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "normalize_url_with_non_default_port",
			input:    "https://example.com:8080/path",
			expected: "https://example.com:8080/path",
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
			<a href="https://other.com/path/two">
				<span>Boot.dev 2</span>
			</a>
		</body>
	</html>
	`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one", "https://other.com/path/two"},
		},
		{
			name:     "default ports and paths",
			inputURL: "https://www.example.com:443",
			inputBody: `
		<html>
			<body>
				<h1>Test URLs with Default Ports and Paths</h1>
				
				<!-- HTTPS URLs with default port 443 -->
				<a href="https://www.example.com:443/page1">HTTPS with explicit 443</a>
				<a href="https://www.example.com:443/api/data">HTTPS API with 443</a>
				<a href="https://www.example.com:443/">HTTPS root with 443</a>
				
				<!-- HTTP URLs with default port 80 -->
				<a href="http://www.example.com:80/page2">HTTP with explicit 80</a>
				<a href="http://www.example.com:80/contact">HTTP contact with 80</a>
				
				<!-- URLs without default ports (should be preserved) -->
				<a href="https://www.example.com:8080/app">Non-default port 8080</a>
				<a href="http://www.example.com:3000/api">Non-default port 3000</a>
				
				<!-- URLs without explicit ports (should work normally) -->
				<a href="https://www.example.com/about">HTTPS without port</a>
				<a href="http://www.example.com/help">HTTP without port</a>
				
				<!-- Mixed case and default ports -->
				<a href="https://WWW.EXAMPLE.COM:443/test">Mixed case with 443</a>
				<a href="http://WWW.EXAMPLE.COM:80/test">Mixed case with 80</a>
				
				<!-- Absolute paths (should resolve relative to base URL) -->
				<a href="/absolute/path">Absolute path</a>
				<a href="/api/users">Absolute API path</a>
				<a href="/static/css/style.css">Absolute static file</a>
				<a href="/">Absolute root</a>
				
				<!-- Relative paths (should resolve relative to base URL) -->
				<a href="relative/page.html">Relative path</a>
				<a href="./current/directory">Current directory</a>
				<a href="../parent/directory">Parent directory</a>
				<a href="../../grandparent/file">Grandparent directory</a>
				<a href="./././nested/current">Nested current directories</a>
				<a href=".././../mixed/path">Mixed relative path</a>
				
				<!-- Paths with query parameters and fragments -->
				<a href="/search?q=test&sort=date">Absolute with query</a>
				<a href="relative/page?param=value#section">Relative with query and fragment</a>
				<a href="https://www.example.com:443/api?key=123&user=john#results">HTTPS with query and fragment</a>
			</body>
		</html>
		`,
			expected: []string{
				"https://www.example.com/page1",
				"https://www.example.com/api/data",
				"https://www.example.com",
				"http://www.example.com/page2",
				"http://www.example.com/contact",
				"https://www.example.com:8080/app",
				"http://www.example.com:3000/api",
				"https://www.example.com/about",
				"http://www.example.com/help",
				"https://www.example.com/test",
				"http://www.example.com/test",
				"https://www.example.com/absolute/path",
				"https://www.example.com/api/users",
				"https://www.example.com/static/css/style.css",
				"https://www.example.com",
				"https://www.example.com/relative/page.html",
				"https://www.example.com/current/directory",
				"https://www.example.com/parent/directory",
				"https://www.example.com/grandparent/file",
				"https://www.example.com/nested/current",
				"https://www.example.com/mixed/path",
				"https://www.example.com/search?q=test&sort=date",
				"https://www.example.com/relative/page?param=value",
				"https://www.example.com/api?key=123&user=john",
			},
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("test %d: expected no error, got %v", i, err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("test %d: expected %v, got %v", i, tc.expected, actual)
			}
		})
	}
}
