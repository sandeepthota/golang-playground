package main

import (
	"fmt"
	"golang.org/x/net/html" //for parsing
	"net/http"              //for crawler
	"strings"               //for has prefix
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
func findLinks(url string) { //[reference :  http://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html]

	baseURL := url[0 : strings.LastIndex(url, "/")+1] //extract base url of web page
	fmt.Println(baseURL)                              // example : http://www.cise.ufl.edu/class/cis4930fa15idm/notes.html to http://www.cise.ufl.edu/class/cis4930fa15idm/
	resp, err := http.Get(url)                        //resp = response, err=error returned if any

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

			fileExt := url[strings.LastIndex(url, ".")+1 : len(url)] //get file extension if any

			//convert link to aboslute link if relative
			if !strings.HasPrefix(url, "http:") {
				//	fmt.Println("Relative link found... Converting to absolute.")
				url = baseURL + url //example notes/dm3part3.pdf to http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm3part2.pdf
			}

			switch {
			case fileExt == "pdf":
				fmt.Println("PDF Found : " + url)

			case fileExt == "ppt":
				fmt.Println("PPT Found : " + url)
			}
			//	fmt.Println(fileExt)
			//	fmt.Println("URL Found : " + url)
		}
	}
}

func main() {
	findLinks("http://www.cise.ufl.edu/class/cis4930fa15idm/notes.html")
	//findLinks("https://www.reddit.com/r/aww")
}
