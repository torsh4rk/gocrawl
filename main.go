package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/html"
)

var url string

func main() {
	// Interrupt the program with Ctrl+C
	handleSignals()
	// Create a new HTTP client with a timeout
	fmt.Printf("\nEnter the URL: ")
	fmt.Scanln(&url)
	if url == "" {
		fmt.Printf("No URL provided. Exiting...\n")
		os.Exit(0)
	}
	fmt.Printf("\n[*] Parsing HTML from %s\n\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	parseHTML(doc)

}

func parseHTML(n *html.Node) {

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			// Check if in the attributes "href", "src","action","data" and "url" there is a link or a path not empty
			if a.Key == "href" || a.Key == "src" || a.Key == "action" || a.Key == "data" || a.Key == "url" {
				// if the attributes is a path not empty, join it with the URL from input and request it via HTTP GET Request and print the status code code not 404
				if a.Val != "" && a.Val != "#" && a.Val != "/" && a.Val != "http" && a.Val != "https" {
					resp, _ := http.Get(url + "/" + a.Val)
					if resp.StatusCode != 404 {
						fmt.Printf("[*] Found: "+"%s%s\n", url, a.Val+" ==> "+resp.Status)
					}
					// if the attributes is a link, only request it via HTTP GET Request and print the status code not 404
				} else if a.Val == "http" || a.Val == "https" {
					resp, _ := http.Get(a.Val)
					if resp.StatusCode != 404 {
						fmt.Printf("[*] Found: "+"%s%s\n", url, a.Val+" ==> "+resp.Status)
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c)
	}

}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\nCtrl+C pressed in Terminal...\nExiting...\n")
		os.Exit(0)
	}()
}
