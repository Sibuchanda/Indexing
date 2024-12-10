package main

import (
	"GoTesting/html_parser" // Replace "GoTesting" with your actual module name
	"fmt"
	"strings"
)

func main() {

	filePath := "./input/geeksforgeeks.html"
	pageURL := "https://www.geeksforgeeks.org"

	metadata, tokens, err := html_parser.ExtractMetaData(filePath, pageURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the metadata and tokens
	fmt.Println("Webpage URL:", metadata.URL)
	fmt.Println("Title:", metadata.Title)
	fmt.Println("Description:", metadata.Description)
	fmt.Println("Keywords:", metadata.Keywords)
	fmt.Println("Tokens:", strings.Join(tokens, ", "))
}
