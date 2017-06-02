package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var url = flag.String("url", "http://usekamba.com", "page url to fetch")

// Get the href attribute from a token in html response Body
// Params: an html Token
// Expected: return true in found or false if not
func getParagraphContent(t html.Token) (paragraph string) {
	return
}

// Extract all paragraphs from a given webpage url
// Params: webpage url
// Expected: return an array of paragraphs
func parseHtml(url string) (paragraphs []string) {
	var isParagraph bool = false
	response, _ := http.Get(url)
	content := response.Body
	defer content.Close()
	z := html.NewTokenizer(content)

	for {
		token := z.Next()
		switch {
		case token == html.ErrorToken:
			return

		case token == html.StartTagToken:
			token := z.Token()
			if token.Data == "p" {
				isParagraph = true
			} else {
				isParagraph = false
			}

		case token == html.TextToken:
			if isParagraph == false {
			} else {
				token := string(z.Text())
				if len(strings.TrimSpace(token)) > 0 {
					paragraphs = append(paragraphs, token)
				}
			}
		}
	}
	return paragraphs
}

func main() {
	flag.Parse()
	paragraphs := parseHtml(*url)
	fmt.Println(len(paragraphs), "paragraphs in", *url, "\n")
	for i, paragraph := range paragraphs {
		fmt.Println(i, paragraph, "\n")
	}
}
