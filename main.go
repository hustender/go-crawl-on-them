package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const baseURL = "https://sdasdas.asd/"

func main() {
	visited := make(map[string]bool)  // To keep track of visited URLs
	dead := make(map[string][]string) // Dead references

	var mu sync.Mutex
	var wg sync.WaitGroup

	crawl(baseURL, baseURL, visited, dead, &mu, &wg)

	wg.Wait()
	printMap(dead, "Site", "Link")
}

// Crawl a URL and recursively fetch links
func crawl(prev string, currentURL string, visited map[string]bool, dead map[string][]string, mu *sync.Mutex, wg *sync.WaitGroup) {
	mu.Lock()
	if visited[currentURL] {
		mu.Unlock()
		return // Skip already visited URLs
	}

	continueCrawl := true

	if !strings.Contains(currentURL, baseURL) {
		continueCrawl = false
	}

	visited[currentURL] = true // Mark URL as visited
	mu.Unlock()

	fmt.Printf("Checking %s for dead links..\n", currentURL)

	wg.Add(1)

	go func() {
		defer wg.Done()

		// Fetch the HTML content of the page
		body := getContent(prev, currentURL, dead)
		if body == "" {
			return
		}

		// Extract hrefs from the page
		hrefs := extractHrefs(body, currentURL)

		// Recursively crawl each extracted link
		for _, href := range hrefs {
			if !continueCrawl {
				return
			}
			crawl(currentURL, href, visited, dead, mu, wg)
		}
	}()
}

// Fetch the content of a URL
func getContent(prev string, pageURL string, dead map[string][]string) string {
	response, err := http.Get(pageURL)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		Dead(dead, prev, pageURL)
		return ""
	}
	if response.StatusCode >= 400 {
		Dead(dead, prev, pageURL)
		return ""
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		Dead(dead, prev, pageURL)
		return ""
	}
	return string(body)
}

func Dead(dead map[string][]string, prev string, pageURL string) {
	dead[prev] = append(dead[prev], pageURL)
}

// Extract hrefs from HTML and resolve them to absolute URLs
func extractHrefs(htmlContent string, base string) []string {
	var hrefs []string

	// Parse the HTML content
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return nil
	}

	// Traverse the HTML nodes and collect href attributes
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					resolvedURL := resolveURL(base, attr.Val)
					if resolvedURL != "" {
						hrefs = append(hrefs, resolvedURL)
					}
					break
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}
	traverse(doc)

	return hrefs
}

// Resolve a relative URL to an absolute URL
func resolveURL(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		fmt.Println("Error parsing base URL:", err)
		return ""
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		fmt.Println("Error parsing href:", err)
		return ""
	}
	return baseURL.ResolveReference(hrefURL).String()
}

func printMap(m map[string][]string, key string, value string) {
	if len(m) == 0 {
		fmt.Println("No dead links found!")
		return
	}

	var maxLenKey int
	for k, _ := range m {
		if len(k) > maxLenKey {
			maxLenKey = len(k)
		}
	}
	fmt.Println(key + ":" + strings.Repeat(" ", max(0, maxLenKey-len(key))) + value + ":")
	for i, j := range m {
		for _, u := range j {
			fmt.Println(i + " " + strings.Repeat(" ", max(0, maxLenKey-len(i))) + u)
		}
	}
}
