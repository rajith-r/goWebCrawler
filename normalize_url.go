package main

import (
	"net/url"
	"path"
	"sort"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	//case normalization
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// fmt.Println("Scheme:", u.Scheme)
	// fmt.Println("Host:", u.Host)
	// fmt.Println("Hostname:", u.Hostname())
	// fmt.Println("Port:", u.Port())

	if u.Port() == "443" && u.Scheme == "https" || u.Port() == "80" && u.Scheme == "http" {
		u.Host = u.Hostname()
		// fmt.Println("New Host:", u.Host)
	}

	// fmt.Println("Path:", u.Path)
	if u.Path != "" {
		u.Path = path.Clean(u.Path) // remove any .. and .
		// fmt.Println("Path after clean:", u.Path)
	}

	if u.Path == "/" {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	// fmt.Println("RawQuery:", u.RawQuery)
	//sort query params
	q := u.Query()
	// fmt.Println("Query:", q)
	// fmt.Println("Length of Query:", len(q))
	keys := make([]string, 0, len(q))
	// fmt.Println("Keys:", keys)
	for k := range q {
		// fmt.Println("Key:", k)
		// fmt.Println("Value:", q.Get(k))
		keys = append(keys, k)
		// fmt.Println("Keys after append:", keys)
	}
	sort.Strings(keys)
	for _, k := range keys {
		q.Set(k, q.Get(k))
		// fmt.Println("Query after set:", q)
	}
	u.RawQuery = q.Encode()

	// 	fmt.Println("RawQuery after sort:", u.RawQuery)

	// fmt.Println("Fragment:", u.Fragment)
	u.Fragment = ""
	// fmt.Println("User:", u.User)
	// fmt.Println("--------------------------------")

	// return u.Scheme + "://" + u.Host, nil
	return u.String(), nil
}
