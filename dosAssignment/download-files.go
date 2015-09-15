//http://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/
//http://stackoverflow.com/questions/8350609/how-do-you-time-a-function-in-go-and-return-its-runtime-in-milliseconds
//https://github.com/bradhe/stopwatch

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	//	"stopwatch" //for timing download function
	"sync" //for wait group
)

var wg sync.WaitGroup

func downloadFromUrl(url string, isMultiThread bool) {
	if isMultiThread {
		defer wg.Done()
	}

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
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
	/*countries := []string{"GB", "FR", "ES", "DE", "CN", "CA", "ID", "US"}
	for i := 0; i < len(countries); i++ {
		url := "http://download.geonames.org/export/dump/" + countries[i] + ".zip"
		downloadFromUrl(url)
	}*/

	listUrl := []string{"http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm1.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm2part1.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm2part2.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm3part1.pdf", "http://www.cise.ufl.edu/class/cis4930fa15idm/notes/dm3part2.pdf"}

	/*for i := 0; i < len(listUrl); i++ {
		wg.Add(0)
		go downloadFromUrl(listUrl[i])
	}*/

	singleThreadDownload(listUrl)

	multiThreadDownload(listUrl)

	wg.Wait()
}
