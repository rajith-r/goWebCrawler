package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	r := strings.NewReader(htmlBody)
	tree, err := html.Parse(r)

	if err != nil {
		return nil, err
	}
	var urls []string
	// fmt.Println("traversing the tree")
	for n := range tree.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			// fmt.Println("n:", n)
			for _, a := range n.Attr {
				// fmt.Println("a:", a)
				// fmt.Println("a.Key:", a.Key)
				// fmt.Println("a.Val:", a.Val)
				// fmt.Println("rawBaseURL:", rawBaseURL)

				if a.Key == "href" {
					url, _ := url.Parse(a.Val)
					if url.Scheme == "" {
						baseURL, err := url.Parse(rawBaseURL)
						// fmt.Println("baseURL:", baseURL)
						if err != nil {
							return nil, err
						}
						url = baseURL.ResolveReference(url)
						// fmt.Println("ResolveReference url:", url)
					}
					// fmt.Println("url:", url)
					if err != nil {
						return nil, err
					}
					u, err := normalizeURL(url.String())
					if err != nil {
						return nil, err
					}
					urls = append(urls, u)
				}
			}
		}
	}
	fmt.Println("urls:", urls)

	return urls, nil
}
