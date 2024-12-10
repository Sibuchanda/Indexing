package html_parser

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery" // For parsing HTML files
	"golang.org/x/text/unicode/norm" // handling Unicode text normalization.
)

type Metadata struct {
	URL         string
	Title       string
	Description string
	Keywords    string
}

// --------- Function for extracting tokens from html body -----------
func Tokenize(text string) []string {
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[a-z]+`)   // Extracting only alphabetic contents from the body
	tokens := re.FindAllString(text, -1) // Finds all substrings in text that match the regular expression ([a-z]+) and stores them as tokens.

	// Function for removing Stop words
	stopWords := map[string]bool{
		"the": true, "is": true, "and": true, "of": true, "a": true, "to": true,
	}

	// This variable stores the filtered tokens (not contains stop words)
	var filteredTokens []string

	// Filtering Out the Stop Words
	for _, token := range tokens {
		if !stopWords[token] {
			filteredTokens = append(filteredTokens, token)
		}
	}
	return filteredTokens
}

// --------------------- Function for extracting metadata from html file --------------------------
func ExtractMetaData(filePath, pageURL string) (Metadata, []string, error) {
	htmlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Metadata{}, nil, fmt.Errorf("error reading file: %v", err)
	}
	htmlContent := string(htmlBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent)) // Parsing the htmlContent string into a goquery document (doc)
	if err != nil {
		return Metadata{}, nil, fmt.Errorf("error loading HTML: %v", err)
	}

	// Extract the title
	title := doc.Find("title").Text()

	// Extract description and keywords from meta tags
	description, _ := doc.Find("meta[name='description']").Attr("content")
	keywords, _ := doc.Find("meta[name='keywords']").Attr("content")

	/*

		// Try extracting the URL from <link> and <base> tags
		var baseURL string

		// Try extracting from <base> tag
		baseURL, _ = doc.Find("base[href]").Attr("href")
		if baseURL == "" {
			// Try extracting URL from the <meta> tag with property 'og:url'
			baseURL, _ = doc.Find("meta[property='og:url']").Attr("content")
		}
		if baseURL == "" {
			// Try extracting URL from the <meta> tag with name 'twitter:url'
			baseURL, _ = doc.Find("meta[name='twitter:url']").Attr("content")
		}
		if baseURL == "" {
			// Try extracting URL from <link> tag with rel 'canonical'
			baseURL, _ = doc.Find("link[rel='canonical']").Attr("href")
		}

		// If no URL found in any of the above methods, try extracting from <a> tags
		if baseURL == "" {
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				// Look for the first valid URL in href
				if strings.HasPrefix(href, "http") {
					baseURL = href
					return
				}
			})
		}

		// If no URL found, return an error
		if baseURL == "" {
			return Metadata{}, nil, fmt.Errorf("URL is not found from the HTML metadata")
		}

		// Clean up the URL by stripping the protocol if present
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return Metadata{}, nil, fmt.Errorf("error parsing URL: %v", err)
		}

		// Strip the protocol (http:// or https://) and just keep the domain
		baseURL = parsedURL.Host

	*/

	// Removing Irrelevant Tags
	doc.Find("script, style, img, iframe").Remove()

	// Extracting all textual content from the <body> tag of the HTML document.
	bodyText := doc.Find("body").Text()

	// Call the tokenize function for extracting tokens from body
	tokens := Tokenize(bodyText)

	metadata := Metadata{
		URL:         pageURL,
		Title:       norm.NFC.String(title),
		Description: norm.NFC.String(description),
		Keywords:    norm.NFC.String(keywords),
	}

	return metadata, tokens, nil
}
