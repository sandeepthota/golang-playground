/*
Some websites cap maximum speed of a connection, hence even if you have higher bandwith, your actual download speed might be lower than your max possible speed. In such ascenario we can use a different thread for each file to be downloaded and imporove our overall speed. In this program, go routines are used to concurrently download files from a list of urls provided.

References :
http://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/

Sample Output : http://showterm.io/ad5ff73e59867afbb888c

*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"sync" //for wait group
	"time"
)

var wg sync.WaitGroup

func downloadFromUrl(url string, isMultiThread bool) { //reference : https://github.com/thbar/golang-playground/blob/master/download-files.go
	if isMultiThread { //check if called by go routine or not
		defer wg.Done()
	}

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

func singleThreadDownload(urlList []string) {
	for i := 0; i < len(urlList); i++ {
		downloadFromUrl(urlList[i], false)
	}
}

func multiThreadDownload(urlList []string) {
	for i := 0; i < len(urlList); i++ {
		wg.Add(1)
		go downloadFromUrl(urlList[i], true)
	}
}

func main() {

	fmt.Println("Starting sample downloads to show speed difference.")

	listUrl := []string{"http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm1.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm2part1.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm2part2.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm3part1.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm3part2.pdf"}

	startTime := time.Now().UTC() //start timer
	singleThreadDownload(listUrl)
	endTime := time.Now().UTC() //stop timer
	var duration1 = endTime.Sub(startTime).Nanoseconds() / 1e6

	startTime2 := time.Now().UTC()
	multiThreadDownload(listUrl)

	wg.Wait() //wait for all go routines to finish

	endTime2 := time.Now().UTC()
	var duration2 = endTime2.Sub(startTime2).Nanoseconds() / 1e6

	fmt.Printf("Time taken in milliseconds for single threaded downloads : %v\n", duration1)
	fmt.Printf("Time taken in milliseconds for multi threaded downloads  : %v\n", duration2)
}
