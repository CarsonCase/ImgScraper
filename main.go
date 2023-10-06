package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	args := os.Args[1:]

	// if it doesn't exist, create img directory
	os.Mkdir("./imgs", 0755)

	c := colly.NewCollector()

	c.OnHTML("img", func(e *colly.HTMLElement) {
		go download(e.Attr("src"))
	})
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(args[0])
}

func download(imageUrl string) {

	segments := strings.Split(imageUrl, "/")
	fileName := segments[len(segments)-1]

	fileNameSegments := strings.Split(fileName, "?")
	fileNameNormal := fileNameSegments[0]

	outputFilePath := "./imgs/" + fileNameNormal // Replace with the desired output file path

	// Make an HTTP GET request to the image URL
	response, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("Error while making the GET request:", err)
		return
	}
	defer response.Body.Close()

	// Check if the response status code is 200 (OK)
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code is not OK:", response.Status)
		return
	}

	// Create or open the output file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error while creating the output file:", err)
		return
	}
	defer outputFile.Close()

	// Copy the image data from the response body to the output file
	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		fmt.Println("Error while copying image data to the output file:", err)
		return
	}

	fmt.Println("Image downloaded successfully to", outputFilePath)
}
