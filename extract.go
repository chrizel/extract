// Package extract can extract parts of HTML pages by using jQuery-like query strings
package extract

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type query struct {
	name    string
	id      string
	classes []string
}

// HTML extracts the HTML code read from r and by applying the query string qs
func HTML(r io.Reader, qs string) (string, error) {
	q := parseQuery(qs)

	doc, err := html.Parse(r)
	if err != nil {
		return "", fmt.Errorf("parse html error: %v", err)
	}

	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if q.matches(n) {
			return nodeToString(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result := f(c)
			if result != "" {
				return result
			}
		}

		return ""
	}

	return f(doc), nil
}

// nodeToString converts a node to an HTML string containing all child data
func nodeToString(node *html.Node) string {
	var b bytes.Buffer
	html.Render(&b, node)

	result := b.String()

	const (
		prefix = "<html><head></head><body>"
		suffix = "</body></html>"
	)

	if strings.HasPrefix(result, prefix) {
		result = result[len(prefix):]
	}
	if strings.HasSuffix(result, suffix) {
		result = result[:len(result)-len(suffix)]
	}

	return result
}

func parseQuery(q string) query {
	result := query{}

	mode := 0 // 0 = name, 1 = id, 2 = class
	var b bytes.Buffer

	next := func(newMode int) {
		s := b.String()
		b.Reset()
		prevMode := mode
		mode = newMode

		if len(s) == 0 {
			return
		}

		switch prevMode {
		case 0:
			result.name = s
		case 1:
			result.id = s
		case 2:
			result.classes = append(result.classes, s)
		}
	}

	for _, c := range q {
		if c == '#' {
			next(1)
			continue
		} else if c == '.' {
			next(2)
			continue
		}

		b.WriteRune(c)
	}
	next(0)

	return result
}

func (q query) matches(n *html.Node) bool {
	return q.matchesName(n) && q.matchesID(n) && q.matchesClasses(n)
}

func (q query) matchesName(n *html.Node) bool {
	if q.name == "" || (n.Type == html.ElementNode && n.Data == q.name) {
		return true
	}

	return false
}

func (q query) matchesID(n *html.Node) bool {
	if q.id == "" {
		return true
	}

	for _, a := range n.Attr {
		if (a.Key == "id" || a.Key == "ID") && a.Val == q.id {
			return true
		}
	}

	return false
}

func (q query) matchesClasses(n *html.Node) bool {
	if len(q.classes) == 0 {
		return true
	}

	for _, a := range n.Attr {
		if a.Key == "class" || a.Key == "CLASS" {
			classes := strings.Split(a.Val, " ")

			matches := true
			for _, class := range q.classes {
				if !stringSliceContains(classes, class) {
					matches = false
				}
			}

			if matches {
				return true
			}
		}
	}

	return false
}

func stringSliceContains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}

	return false
}
