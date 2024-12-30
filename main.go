package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	// Execute
	var rootCmd = &cobra.Command{
		Use:   "crawl",
		Short: "go-crawl-on-them is a basic web-crawler designed to find dead links inside a website",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Usage: crawl <url>")
				return
			}
			baseURL = args[0]
			fmt.Printf("Crawling: '%s'\n", baseURL)
			run()
		},
	}

	// Error Handling
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var baseURL string

func run() {
	visited := make(map[string]bool)  // To keep track of visited URLs
	dead := make(map[string][]string) // Dead references

	// goroutines
	var mu sync.Mutex
	var wg sync.WaitGroup

	crawl(baseURL, baseURL, visited, dead, &mu, &wg) // Begin crawling
	wg.Wait()                                        // Let goroutines finish

	printMap(dead, "Site", "Link") // Print results
}

// Crawl a URL and recursively fetch links
func crawl(prev string, currentURL string, visited map[string]bool, dead map[string][]string, mu *sync.Mutex, wg *sync.WaitGroup) {
	mu.Lock()
	if visited[currentURL] {
		mu.Unlock()
		return // Skip already visited URLs
	}

	visited[currentURL] = true // Mark URL as visited
	mu.Unlock()

	fmt.Printf("Checking '%s' for dead links..\n", currentURL)

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
			if !strings.Contains(currentURL, baseURL) {
				return // We still want a reference to a different website to be saved but not to be crawled
			}
			crawl(currentURL, href, visited, dead, mu, wg) // Crawl this website too!
		}
	}()
}

// Fetch the content of a URL
func getContent(prev string, pageURL string, dead map[string][]string) string {
	response, err := request(pageURL) // Request to the url

	// Error handling
	if err != nil {
		fmt.Printf("Error fetching URL '%s': %s\n", pageURL, err)
		deadLink(dead, prev, pageURL) // Dead link
		return ""
	}
	if response.StatusCode >= 400 {
		deadLink(dead, prev, pageURL)
		return ""
	}

	// Read response body
	body, err := io.ReadAll(response.Body)
	// Error handling
	if err != nil {
		fmt.Println("Error reading response body:", err)
		deadLink(dead, prev, pageURL)
		return ""
	}
	return string(body) // Return body content
}

// Request to the url
func request(pageURL string) (*http.Response, error) {
	// Response
	type ResponseResult struct {
		Response *http.Response
		Err      error
	}

	// Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan ResponseResult, 1)

	// Handle possible timeout
	go func() {
		response, err := http.Get(pageURL)
		resultChan <- ResponseResult{response, err}
	}()
	select {
	case <-ctx.Done():
		return nil, http.ErrHandlerTimeout
	case result := <-resultChan:
		return result.Response, result.Err
	}
}

// Marks a link as dead
func deadLink(dead map[string][]string, prev string, pageURL string) {
	dead[prev] = append(dead[prev], pageURL) // Append dead link
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

// Prints the results
func printMap(m map[string][]string, key string, value string) {
	if len(m) == 0 {
		fmt.Println("No dead links found!")
		return
	}

	// Pretty printing
	var maxLenKey int
	for k := range m {
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
