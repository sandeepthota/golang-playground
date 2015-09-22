// http://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html

package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

//get href attribute from token
func getHref(t html.Token) (notPresent bool, href string) {
	//iterate  over token's attributes till you find an href
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			notPresent = false
		}
	}
	return
}

//find links from webpage and display
func findLinks(url string) {
	resp, err := http.Get(url) //resp = response, err=error returned if any

	if err != nil {
		fmt.Println("Error in crawling url : " + url)
		return
	}

	b := resp.Body
	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			//reached end of document
			return
		case tt == html.StartTagToken:
			t := z.Token()

			//check if token is a hyper link tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			//extract the linked url in hyperlink tag if present
			notPresent, url := getHref(t)
			if notPresent {
				continue
			}

			fmt.Println("URL Found : " + url)
		}
	}
}

func main() {
	findLinks("http://www.cise.ufl.edu/class/cis4930fa15idm/notes.html")
	//findLinks("https://www.reddit.com/r/aww")
}
