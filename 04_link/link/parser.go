package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a HTML link
type Link struct {
	Href string
	Text string
}

// Parse returns all the links found in the HTML from the given Reader
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var links []Link
	var text []string

	// Adapted from depth-first order processing example in https://godoc.org/golang.org/x/net/html
	var searchLinks func(*html.Node, bool)
	searchLinks = func(n *html.Node, grabText bool) {
		isLink := n.Type == html.ElementNode && n.Data == "a"

		if n.Type == html.TextNode && grabText {
			text = append(text, strings.TrimSpace(n.Data))
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			searchLinks(c, grabText || isLink)
		}

		if isLink {
			href := attributeValue(n, "href")
			if href != nil {
				links = append(links, Link{Href: *href, Text: strings.TrimSpace(strings.Join(text, " "))})
			}
			text = nil
		}
	}
	searchLinks(doc, false)

	return links, nil
}

func attributeValue(n *html.Node, key string) *string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return &attr.Val
		}
	}

	return nil
}
