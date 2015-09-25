/*
File Downloader : Given link to a webpage and file type, this program will download all files of those extension(if any) present on the webpage. For each link found containing file of desired extension, a go routine is fired. This will ensure maximum throughput since some of the files found might be present on a slower server while others might be present on a faster server.

Running Instructions :
Run go get golang.org/x/net/html on command prompt first to get html package. The program takes url and file extension as command line arguments. So, to get all jpg files from https://www.reddit.com/r/pics you would type in go run fileDownloader.go https://www.reddit.com/r/pics jpg

Sample Output : http://showterm.io/500cb9fee78a64d7f09e4

*/

package main

import (
	"fmt"
	"golang.org/x/net/html" //for parsing.
	"io"
	"net/http" //for crawler
	"os"       //creating file
	"strings"  //for has prefix
	"sync"     //for wait group
)

var wg sync.WaitGroup

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

//find links to files which have specified extension
func findFiles(url string, extension string) { //[reference :  http://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html]

	baseURL := url[0 : strings.LastIndex(url, "/")+1] //extract base url of web page
	//	fmt.Println(baseURL)                              // example : http://www.cise.ufl.edu/class/cis4930fa15idm/notes.html to http://www.cise.ufl.edu/class/cis4930fa15idm/
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

			fileExt := url[strings.LastIndex(url, ".")+1 : len(url)] //get file extension if any

			//convert link to aboslute link if relative
			if !strings.HasPrefix(url, "http:") {
				//	fmt.Println("Relative link found... Converting to absolute.")
				url = baseURL + url //example notes/dm3part3.pdf to http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm3part2.pdf
			}

			if fileExt == extension {
				fmt.Println("File of ." + fileExt + " extension found : " + url)
				wg.Add(1)
				go downloadFromUrl(url)
			}
		}
	}
}

func downloadFromUrl(url string) { //reference : https://github.com/thbar/golang-playground/blob/master/download-files.go

	defer wg.Done()

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func main() {
	findFiles(os.Args[1], os.Args[2])
	wg.Wait() //wait for all go routines to finish
}
